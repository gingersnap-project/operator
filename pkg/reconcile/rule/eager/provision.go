package eager

import (
	"fmt"
	"strings"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/meta"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	apicorev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

func LoadCache(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cacheRef := r.Spec.Cache
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

func DBSyncer(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cache := ctx.Cache
	// TODO use globally unique naming convention as Rules with the same names can be created across multiple namespaces
	name := r.Name
	labels := meta.GingersnapLabels("db-syncer", meta.ComponentDBSyncer, r.Name)

	secret := &apicorev1.Secret{}
	if err := ctx.Client().Load(cache.ConfigurationSecret(), secret); err != nil {
		ctx.Requeue(fmt.Errorf("unable to load cache configuratioon secret: %w", err))
	}

	// TODO how to provide properties securely? application.properties mounted secret?
	host := secret.Data["host"]
	port := secret.Data["port"]
	user := secret.Data["username"]
	password := secret.Data["password"]
	javaOpts := []string{
		fmt.Sprintf("-Dgingersnap.region.us-east.backend.uri=hotrod://%s:%s@%s:%s?sasl_mechanism=SCRAM-SHA-512", user, password, host, port),
		"-Dgingersnap.region.us-east.database.hostname=mysql.mysql.svc.cluster.local",
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
							WithEnv(
								corev1.EnvVar().WithName("JAVA_OPTS").WithValue(strings.Join(javaOpts, " ")),
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
