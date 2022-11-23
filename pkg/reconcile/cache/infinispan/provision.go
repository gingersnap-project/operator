package infinispan

import (
	"fmt"
	"strconv"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	monitoringv1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/monitoring/v1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache/context"
	"github.com/gingersnap-project/operator/pkg/reconcile/meta"
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

func WatchServiceAccount(c *v1alpha1.Cache, ctx *context.Context) {
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
				WithVerbs("watch"),
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

func Service(c *v1alpha1.Cache, ctx *context.Context) {
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

func ConfigurationSecret(c *v1alpha1.Cache, ctx *context.Context) {
	cacheService := c.CacheService()
	secretName := cacheService.ConfigurationSecret()
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

	labels := resourceLabels(c)
	secret := corev1.Secret(secretName, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithStringData(
			map[string]string{
				"type":     "gingersnap",
				"provider": "gingersnap",
				"host":     cacheService.SvcName(),
				"port":     strconv.Itoa(8080),
			},
		).
		WithType("servicebinding.io/gingersnap")

	if err := ctx.Client().Apply(secret); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan configuration secret: %w", err))
		return
	}

	c.Status.ServiceBinding = &v1alpha1.ServiceBinding{
		Name: secretName,
	}
	if err := ctx.Client().UpdateStatus(c); err != nil {
		ctx.Requeue(fmt.Errorf("unable to add ServiceBinding to Cache Status CR: %w", err))
	}
}

func DaemonSet(c *v1alpha1.Cache, ctx *context.Context) {
	labels := resourceLabels(c)
	ds := appsv1.
		DaemonSet(c.Name, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DaemonSetSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(labels),
			).
			WithTemplate(corev1.PodTemplateSpec().
				WithName(sidecarContainerName).
				WithLabels(labels).
				WithSpec(corev1.PodSpec().
					WithServiceAccountName(c.Name).
					WithContainers(
						corev1.Container().
							WithName(sidecarContainerName).
							WithImage("quay.io/gingersnap/cache-manager").
							WithPorts(
								corev1.ContainerPort().WithContainerPort(8080),
								corev1.ContainerPort().WithContainerPort(11222),
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
				),
			),
		)
	if err := ctx.Client().Apply(ds); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan DaemonSet: %w", err))
	}
}

func ServiceMonitor(c *v1alpha1.Cache, ctx *context.Context) {
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
										Name: c.CacheService().ConfigurationSecret(),
									},
									Key: "password",
								},
							).
							WithUsername(
								apicorev1.SecretKeySelector{
									LocalObjectReference: apicorev1.LocalObjectReference{
										Name: c.CacheService().ConfigurationSecret(),
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
