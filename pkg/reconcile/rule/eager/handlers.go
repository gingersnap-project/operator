package eager

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	bindingv1 "github.com/gingersnap-project/operator/pkg/applyconfigurations/binding/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/images"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile/meta"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	apiappsv1 "k8s.io/api/apps/v1"
	apicorev1 "k8s.io/api/core/v1"
	apirbacv1 "k8s.io/api/rbac/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	rbacv1 "k8s.io/client-go/applyconfigurations/rbac/v1"
	"k8s.io/utils/pointer"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func LoadCache(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cacheRef := r.CacheService()
	cache := &v1alpha1.Cache{}
	err := ctx.Client().
		WithNamespace(cacheRef.Namespace).
		Load(cacheRef.Name, cache)

	if err != nil {
		msg := fmt.Sprintf("unable to load Cache CR '%s'", cacheRef)
		r.SetCondition(
			v1alpha1.EagerCacheRuleCondition{
				Type:    v1alpha1.EagerCacheRuleConditionReady,
				Status:  apimetav1.ConditionFalse,
				Message: msg,
			},
		)
		if err := ctx.Client().UpdateStatus(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to update Ready condition on LoadCache failure: %w", err))
			return
		}
		ctx.Requeue(fmt.Errorf("%s: %w", msg, err))
	}
	ctx.Cache = cache
}

func ApplyDBServiceBinding(_ *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cache := ctx.Cache
	labels := meta.GingersnapLabels("db-syncer", meta.ComponentDBSyncer, cache.Name)

	var serviceRef *bindingv1.ServiceApplyConfiguration
	ds := cache.Spec.DataSource
	if service := ds.ServiceProviderRef; service != nil {
		// Webhook validation ensures that parsing never fails
		groupVersion, _ := schema.ParseGroupVersion(service.ApiVersion)
		serviceRef = bindingv1.Service().
			WithGroup(groupVersion.Group).
			WithVersion(groupVersion.Version).
			WithKind(service.Kind).
			WithName(service.Name)
	} else {
		serviceRef = bindingv1.Service().
			WithGroup(apicorev1.GroupName).
			WithVersion(apicorev1.SchemeGroupVersion.Version).
			WithKind("Secret").
			WithName(ds.SecretRef.Name)
	}

	sb := bindingv1.ServiceBinding(cache.CacheService().DBSyncerDataServiceBinding(), cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(
			bindingv1.ServiceBindingSpec().
				WithServices(serviceRef).
				WithApplication(
					bindingv1.Application().
						WithGroup(apiappsv1.GroupName).
						WithVersion(apiappsv1.SchemeGroupVersion.Version).
						WithKind("Deployment").
						WithLabelSelector(
							apimetav1.LabelSelector{
								MatchLabels: labels,
							},
						),
				),
		)

	if err := ctx.Client().Apply(sb); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply DB ServiceBinding: %w", err))
	}
}

func ApplyCacheServiceBinding(_ *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cache := ctx.Cache
	labels := meta.GingersnapLabels("db-syncer", meta.ComponentDBSyncer, cache.Name)

	sb := bindingv1.ServiceBinding(cache.CacheService().DBSyncerCacheServiceBinding(), cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(
			bindingv1.ServiceBindingSpec().
				WithServices(
					bindingv1.Service().
						WithGroup(apicorev1.GroupName).
						WithVersion(apicorev1.SchemeGroupVersion.Version).
						WithKind("Secret").
						WithName(cache.CacheService().DBSyncerCacheServiceBindingSecret()),
				).
				WithApplication(
					bindingv1.Application().
						WithGroup(apiappsv1.GroupName).
						WithVersion(apiappsv1.SchemeGroupVersion.Version).
						WithKind("Deployment").
						WithLabelSelector(
							apimetav1.LabelSelector{
								MatchLabels: labels,
							},
						),
				),
		)

	if err := ctx.Client().Apply(sb); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Cache ServiceBinding: %w", err))
	}
}

func ApplyDBSyncer(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cache := ctx.Cache
	labels := meta.GingersnapLabels("db-syncer", meta.ComponentDBSyncer, cache.Name)

	cacheService := cache.CacheService()
	name := cacheService.DBSyncerName()
	deployment := appsv1.Deployment(name, cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DeploymentSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(labels),
			).
			WithTemplate(corev1.PodTemplateSpec().
				WithName("db-syncer").
				WithLabels(labels).
				WithSpec(corev1.PodSpec().
					WithServiceAccountName(cache.Name).
					WithContainers(
						corev1.Container().
							WithName("db-syncer").
							WithImage(images.DBSyncer).
							WithEnv(
								corev1.EnvVar().WithName("GINGERSNAP_DYNAMIC_MEMBERSHIP").WithValue("true"),
								corev1.EnvVar().WithName("GINGERSNAP_K8S_NAMESPACE").WithValue(cacheService.Namespace),
								corev1.EnvVar().WithName("GINGERSNAP_K8S_RULE_CONFIG_MAP").WithValue(cacheService.EagerCacheConfigMap()),
								corev1.EnvVar().WithName("QUARKUS_LOG_CATEGORY__IO_QUARKUS_KUBERNETES_SERVICE_BINDING__LEVEL").WithValue("DEBUG"),
							).
							WithResources(
								corev1.ResourceRequirements().
									WithLimits(cache.DBSyncerLimits()).
									WithRequests(cache.DBSyncerRequests()),
							).
							WithVolumeMounts(
								corev1.VolumeMount().WithName("eager-rules").WithMountPath("/rules/eager").WithReadOnly(true),
							),
					).
					WithVolumes(
						corev1.Volume().
							WithName("eager-rules").
							WithConfigMap(
								corev1.ConfigMapVolumeSource().WithName(r.ConfigMap()).WithOptional(true),
							),
					),
				),
			),
		)

	if err := ctx.Client().Apply(deployment); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply DB-Syncer Deployment: %w", err))
	}
}

// TODO configure credentials via Secret volume mount
// TODO allow resource req/limits to be configured via Cache CR
// TODO only provision if QueryEnabled?
func ApplySearchIndex(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	if !r.Query() {
		return
	}
	cache := ctx.Cache
	labels := meta.GingersnapLabels(meta.ComponentSearchIndex, meta.ComponentSearchIndex, cache.Name)

	name := cache.CacheService().SearchIndexName()
	deployment := appsv1.Deployment(name, cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DeploymentSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(labels),
			).
			WithTemplate(corev1.PodTemplateSpec().
				WithName(meta.ComponentSearchIndex).
				WithLabels(labels).
				WithSpec(corev1.PodSpec().
					WithContainers(
						corev1.Container().
							WithName(meta.ComponentSearchIndex).
							WithImage(images.Index).
							WithEnv(
								corev1.EnvVar().WithName("discovery.type").WithValue("single-node"),
								corev1.EnvVar().WithName("plugins.security.ssl.http.enabled").WithValue("false"),
							),
					),
				),
			),
		)

	if ctx.Openshift() {
		// We must create a ServiceAccount and RoleBinding so that we can allow the pod to execute with the "anyuid"
		// Security Context Constraint
		// This is necessary as the opensearch image needs to be executed with uid and gid = 1000
		// https://github.com/opensearch-project/opensearch-devops/issues/97
		serviceAccount := corev1.ServiceAccount(name, cache.Namespace)
		if err := ctx.Client().Apply(serviceAccount); err != nil {
			ctx.Requeue(fmt.Errorf("unable to apply Index ServiceAccount: %w", err))
			return
		}

		roleBinding := rbacv1.RoleBinding(name, cache.Namespace).
			WithRoleRef(
				rbacv1.RoleRef().
					WithAPIGroup("rbac.authorization.k8s.io").
					WithKind("ClusterRole").
					WithName("system:openshift:scc:anyuid"),
			).
			WithSubjects(
				rbacv1.Subject().
					WithKind(apirbacv1.ServiceAccountKind).
					WithName(name).
					WithNamespace(cache.Namespace),
			)

		if err := ctx.Client().Apply(roleBinding); err != nil {
			ctx.Requeue(fmt.Errorf("unable to apply Index RoleBinding: %w", err))
			return
		}

		deployment.Spec.Template.Spec.ServiceAccountName = serviceAccount.Name
		deployment.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContextApplyConfiguration{
			RunAsGroup: pointer.Int64(1000),
			RunAsUser:  pointer.Int64(1000),
		}
	}

	if err := ctx.Client().Apply(deployment); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Index Deployment: %w", err))
	}
}

func ApplySearchIndexService(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cache := ctx.Cache
	labels := meta.GingersnapLabels(meta.ComponentSearchIndex, meta.ComponentSearchIndex, cache.Name)
	service := corev1.
		Service(cache.CacheService().SearchIndexSvcName(), cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithSpec(
			corev1.ServiceSpec().
				WithType(apicorev1.ServiceTypeClusterIP).
				WithSelector(labels).
				WithPorts(
					corev1.ServicePort().WithName("opensearch").WithPort(9200),
				),
		)

	if err := ctx.Client().Apply(service); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply SearchIndexService: %w", err))
	}
}

func RemoveDBSyncer(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cacheService := r.CacheService()
	labels := cacheService.LabelSelector()
	eagerCaches := &v1alpha1.EagerCacheRuleList{}
	if err := ctx.Client().List(labels, eagerCaches, client.ClusterScoped); err != nil {
		ctx.Requeue(fmt.Errorf("unable to list all EagerCacheRules to determine db-syncer lifecycle: %w", err))
		return
	}

	if len(eagerCaches.Items) == 1 && eagerCaches.Items[0].UID == r.UID {
		// Remove the db-syncer deployment as no other dependent EagerCacheRules exist
		if err := ctx.Client().Delete(cacheService.DBSyncerName(), &apiappsv1.Deployment{}); runtimeClient.IgnoreNotFound(err) != nil {
			ctx.Requeue(fmt.Errorf("unable to remove db-syncer: %w", err))
		}
	}
}

func RemoveSearchIndex(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cacheService := r.CacheService()
	labels := cacheService.LabelSelector()
	eagerCaches := &v1alpha1.EagerCacheRuleList{}
	if err := ctx.Client().List(labels, eagerCaches, client.ClusterScoped); err != nil {
		ctx.Requeue(fmt.Errorf("unable to list all EagerCacheRules to determine Index lifecycle: %w", err))
		return
	}

	if len(eagerCaches.Items) == 1 && eagerCaches.Items[0].UID == r.UID {
		// Remove the index deployment as no other dependent EagerCacheRules exist
		if err := ctx.Client().Delete(cacheService.SearchIndexName(), &apiappsv1.Deployment{}); runtimeClient.IgnoreNotFound(err) != nil {
			ctx.Requeue(fmt.Errorf("unable to remove Index Deployment: %w", err))
		}

		// Remove the index service
		if err := ctx.Client().Delete(cacheService.SearchIndexSvcName(), &apicorev1.Service{}); runtimeClient.IgnoreNotFound(err) != nil {
			ctx.Requeue(fmt.Errorf("unable to remove Index Service: %w", err))
		}
	}
}
