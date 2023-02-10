package cache

import (
	"fmt"
	"strconv"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	monitoringv1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/monitoring/v1"
	bindingv1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/servicebinding/v1beta1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/meta"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	apirbacv1 "k8s.io/api/rbac/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	rbacv1 "k8s.io/client-go/applyconfigurations/rbac/v1"
)

const (
	sidecarContainerName = "cache-manager"
)

func resourceLabels(c *v1alpha1.Cache) map[string]string {
	return meta.GingersnapLabels("infinispan", meta.ComponentCache, c.Name)
}

func WatchServiceAccount(c *v1alpha1.Cache, ctx *Context) {
	serviceAccount := corev1.ServiceAccount(c.Name, c.Namespace).
		WithOwnerReferences(ctx.Client().OwnerReference())

	if err := ctx.Client().Apply(serviceAccount); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply ServiceAccount: %w", err))
		return
	}

	role := rbacv1.Role(c.Name, c.Namespace).
		WithRules(
			rbacv1.PolicyRule().
				WithAPIGroups("").
				WithResources("configmaps").
				WithVerbs("get", "watch", "list"),
		).
		WithOwnerReferences(ctx.Client().OwnerReference())

	if err := ctx.Client().Apply(role); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Role: %w", err))
		return
	}

	roleBinding := rbacv1.RoleBinding(c.Name, c.Namespace).
		WithRoleRef(
			rbacv1.RoleRef().
				WithAPIGroup("rbac.authorization.k8s.io").
				WithKind("Role").
				WithName(c.Name),
		).
		WithSubjects(
			rbacv1.Subject().
				WithKind(apirbacv1.ServiceAccountKind).
				WithName(c.Name).
				WithNamespace(c.Namespace),
		).
		WithOwnerReferences(ctx.Client().OwnerReference())

	if err := ctx.Client().Apply(roleBinding); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply RoleBinding: %w", err))
		return
	}
}

func Service(c *v1alpha1.Cache, ctx *Context) {
	labels := resourceLabels(c)
	service := corev1.
		Service(c.Name, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithSpec(
			corev1.ServiceSpec().
				WithClusterIP(apicorev1.ClusterIPNone).
				WithInternalTrafficPolicy(apicorev1.ServiceInternalTrafficPolicyLocal).
				WithType(apicorev1.ServiceTypeClusterIP).
				WithSelector(labels).
				WithPorts(
					corev1.ServicePort().WithName("hotrod").WithPort(11222),
					corev1.ServicePort().WithName("rest").WithPort(8080),
				),
		)

	if err := ctx.Client().Apply(service); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan Service: %w", err))
	}
}

func ApplyDataSourceServiceBinding(cache *v1alpha1.Cache, ctx *Context) {
	labels := resourceLabels(cache)

	var serviceRef *bindingv1.ServiceBindingServiceReferenceApplyConfiguration
	ds := cache.Spec.DataSource
	if service := ds.ServiceProviderRef; service != nil {
		serviceRef = bindingv1.ServiceBindingServiceReference().
			WithAPIVersion(service.ApiVersion).
			WithKind(service.Kind).
			WithName(service.Name)
	} else {
		serviceRef = bindingv1.ServiceBindingServiceReference().
			WithAPIVersion(apicorev1.SchemeGroupVersion.String()).
			WithKind("Secret").
			WithName(ds.SecretRef.Name)
	}

	var workloadKind string
	if cache.Local() {
		workloadKind = "DaemonSet"
	} else {
		workloadKind = "Deployment"
	}

	sb := bindingv1.ServiceBinding(cache.CacheService().DataSourceServiceBinding(), cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(
			bindingv1.ServiceBindingSpec().
				WithService(serviceRef).
				WithType(ds.DbType.ServiceBinding()).
				WithWorkload(
					bindingv1.ServiceBindingWorkloadReference().
						WithAPIVersion(apiappsv1.SchemeGroupVersion.String()).
						WithKind(workloadKind).
						WithName(cache.Name),
				),
		)

	if err := ctx.Client().Apply(sb); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Cache ServiceBinding: %w", err))
	}
}

func UserServiceBindingSecret(c *v1alpha1.Cache, ctx *Context) {
	cacheService := c.CacheService()
	secretName := cacheService.UserServiceBindingSecret()
	// TODO reinstate once Authentication has been added to cache-manager
	//existingSecret := &apicorev1.Secret{}
	//var password string
	//if err := ctx.Client().Load(secretName, existingSecret); err != nil {
	//	if !errors.IsNotFound(err) {
	//		ctx.Requeue(fmt.Errorf("unable to retrieve existing configuration secret: %w", err))
	//		return
	//	}
	//
	//	password, err = passwords.Generate(16)
	//	if err != nil {
	//		ctx.Requeue(fmt.Errorf("unable to generate password: %w", err))
	//		return
	//	}
	//} else {
	//	// Configuration secret already exists, so make sure we apply the existing password
	//	password = string(existingSecret.Data["password"])
	//}

	// Initialize the ctx ServiceBinding so that we can use the values when creating the DaemonSet
	secret := serviceBindingSecret(secretName, 8080, c, ctx)

	if err := ctx.Client().Apply(secret); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply user ServiceBinding secret: %w", err))
		return
	}

	c.Status.ServiceBinding = &v1alpha1.ServiceBinding{
		Name: secretName,
	}
	if err := ctx.Client().UpdateStatus(c); err != nil {
		ctx.Requeue(fmt.Errorf("unable to add ServiceBinding to Cache Status CR: %w", err))
	}
}

func DBSyncerCacheServiceBindingSecret(c *v1alpha1.Cache, ctx *Context) {
	// TODO add authentication details once implemented in cache-manager
	secretName := c.CacheService().DBSyncerCacheServiceBindingSecret()
	secret := serviceBindingSecret(secretName, 11222, c, ctx)

	if err := ctx.Client().Apply(secret); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply internal ServiceBinding secret: %w", err))
		return
	}
}

func serviceBindingSecret(name string, port int, c *v1alpha1.Cache, ctx *Context) *corev1.SecretApplyConfiguration {
	labels := resourceLabels(c)
	return corev1.Secret(name, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithStringData(
			map[string]string{
				"type":     "gingersnap",
				"provider": "gingersnap",
				"host":     c.CacheService().SvcName(),
				"port":     strconv.Itoa(port),
			},
		).
		WithType("servicebinding.io/gingersnap")
}

func Deployment(c *v1alpha1.Cache, ctx *Context) {
	labels := resourceLabels(c)
	deployment := appsv1.
		Deployment(c.Name, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DeploymentSpec().
			WithReplicas(c.Spec.Deployment.Replicas).
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(labels),
			).
			WithTemplate(podTemplateSpec(c)),
		)
	if err := ctx.Client().Apply(deployment); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan DaemonSet: %w", err))
	}
}

func DaemonSet(c *v1alpha1.Cache, ctx *Context) {
	labels := resourceLabels(c)
	ds := appsv1.
		DaemonSet(c.Name, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DaemonSetSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(labels),
			).
			WithTemplate(podTemplateSpec(c)),
		)
	if err := ctx.Client().Apply(ds); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan DaemonSet: %w", err))
	}
}

func podTemplateSpec(c *v1alpha1.Cache) *corev1.PodTemplateSpecApplyConfiguration {
	return corev1.PodTemplateSpec().
		WithName(sidecarContainerName).
		WithLabels(resourceLabels(c)).
		WithSpec(corev1.PodSpec().
			WithServiceAccountName(c.Name).
			WithContainers(
				corev1.Container().
					WithName(sidecarContainerName).
					WithImage(c.CacheManagerImage()).
					WithEnv(
						corev1.EnvVar().WithName("GINGERSNAP_K8S_EAGER_CONFIG_MAP").WithValue(c.CacheService().EagerCacheConfigMap()),
						corev1.EnvVar().WithName("GINGERSNAP_K8S_LAZY_CONFIG_MAP").WithValue(c.CacheService().LazyCacheConfigMap()),
						corev1.EnvVar().WithName("GINGERSNAP_K8S_NAMESPACE").WithValue(c.Namespace),
						corev1.EnvVar().WithName("QUARKUS_LOG_CATEGORY__IO_QUARKUS_KUBERNETES_SERVICE_BINDING__LEVEL").WithValue("DEBUG"),
					).
					WithPorts(
						corev1.ContainerPort().WithContainerPort(8080),
						corev1.ContainerPort().WithContainerPort(11222),
					).
					WithResources(
						corev1.ResourceRequirements().
							WithLimits(c.DeploymentLimits()).
							WithRequests(c.DeploymentRequests()),
					).
					WithLivenessProbe(
						httpProbe("live", 5, 0, 10, 1, 80, 8080),
					).
					WithReadinessProbe(
						httpProbe("ready", 5, 0, 10, 1, 80, 8080),
					).
					WithStartupProbe(
						httpProbe("started", 600, 1, 1, 1, 80, 8080),
					).
					WithVolumeMounts(
						corev1.VolumeMount().WithName("lazy-rules").WithMountPath("/rules/lazy").WithReadOnly(true),
					),
			).
			WithVolumes(
				corev1.Volume().
					WithName("lazy-rules").
					WithConfigMap(
						corev1.ConfigMapVolumeSource().WithName(c.CacheService().LazyCacheConfigMap()).WithOptional(true),
					),
			),
		)
}

func ServiceMonitor(c *v1alpha1.Cache, ctx *Context) {
	if !ctx.IsTypeSupported(reconcile.ServiceMonitorGVK) {
		return
	}

	labels := resourceLabels(c)
	serviceMonitor := monitoringv1.
		ServiceMonitor(c.Name, c.Namespace).
		WithLabels(labels).
		WithSpec(monitoringv1.ServiceMonitorSpec().
			WithEndpoints(
				monitoringv1.Endpoint().
					WithBasicAuth(
						monitoringv1.BasicAuth().
							WithPassword(
								apicorev1.SecretKeySelector{
									LocalObjectReference: apicorev1.LocalObjectReference{
										Name: c.CacheService().UserServiceBindingSecret(),
									},
									Key: "password",
								},
							).
							WithUsername(
								apicorev1.SecretKeySelector{
									LocalObjectReference: apicorev1.LocalObjectReference{
										Name: c.CacheService().UserServiceBindingSecret(),
									},
									Key: "username",
								},
							),
					).
					WithHonorLabels(true).
					WithInterval("30s").
					WithPath("/metrics").
					WithPort("infinispan").
					WithScheme("http").
					WithScrapeTimeout("10s"),
			).
			WithNamespaceSelector(
				monitoringv1.NamespaceSelector().WithMatchNames(c.Namespace),
			).
			WithSelector(
				apimetav1.LabelSelector{
					MatchLabels: labels,
				},
			),
		)
	if err := ctx.Client().Apply(serviceMonitor); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan ServiceMonitor: %w", err))
	}
}

func httpProbe(endpoint string, failureThreshold, initialDelay, period, successThreshold, timeout int32, port int) *corev1.ProbeApplyConfiguration {
	return corev1.Probe().
		WithHTTPGet(
			corev1.HTTPGetAction().
				WithScheme(apicorev1.URISchemeHTTP).
				WithPath("q/health/" + endpoint).
				WithPort(intstr.FromInt(port)),
		).
		WithFailureThreshold(failureThreshold).
		WithInitialDelaySeconds(initialDelay).
		WithPeriodSeconds(period).
		WithSuccessThreshold(successThreshold).
		WithTimeoutSeconds(timeout)
}
