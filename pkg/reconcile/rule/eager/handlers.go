package eager

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile/meta"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	apiappsv1 "k8s.io/api/apps/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func LoadCache(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cacheRef := r.CacheService()
	cache := &v1alpha1.Cache{}
	err := ctx.Client().
		WithNamespace(cacheRef.Namespace).
		Load(cacheRef.Name, cache)

	if err != nil {
		// TODO set status !Ready condition
		ctx.Requeue(fmt.Errorf("unable to load Cache CR '%s': %w", cacheRef, err))
	}
	ctx.Cache = cache
}

func ApplyDBSyncer(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cache := ctx.Cache
	labels := meta.GingersnapLabels("db-syncer", meta.ComponentDBSyncer, cache.Name)

	appProps := fmt.Sprintf(`
		gingersnap.rule.us-east.backend.uri=hotrod://%s:%d
		gingersnap.rule.us-east.database.hostname=mysql.mysql.svc.cluster.local
`,
		cache.CacheService().SvcName(),
		11222,
	)

	name := cache.CacheService().DBSyncerName()
	propsSecret := corev1.Secret(name, cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithStringData(map[string]string{
			"application.properties": appProps,
		})

	if err := ctx.Client().Apply(propsSecret); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply DB-Syncer config Secret: %w", err))
		return
	}

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
							WithImage("quay.io/gingersnap/db-syncer").
							WithResources(
								corev1.ResourceRequirements().
									WithLimits(cache.DBSyncerLimits()).
									WithRequests(cache.DBSyncerRequests()),
							).
							WithVolumeMounts(
								corev1.VolumeMount().WithName("config").WithMountPath("/deployments/config").WithReadOnly(true),
								corev1.VolumeMount().WithName("eager-rules").WithMountPath("/rules/eager").WithReadOnly(true),
							),
					).
					WithVolumes(
						corev1.Volume().
							WithName("config").
							WithSecret(
								corev1.SecretVolumeSource().WithSecretName(*propsSecret.Name),
							),
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
