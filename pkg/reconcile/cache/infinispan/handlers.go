package infinispan

import (
	"fmt"

	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/infinispan/configuration"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const (
	containerName = "infinispan"
)

var (
	selectorLabels = map[string]string{
		"app": "infinispan",
	}
)

func Service(c *v1alpha1.Cache, ctx reconcile.Context) {
	service := corev1.
		Service(c.Name, c.Namespace).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithSpec(
			corev1.ServiceSpec().
				WithClusterIP(apicorev1.ClusterIPNone).
				WithType(apicorev1.ServiceTypeClusterIP).
				WithSelector(selectorLabels).
				WithPorts(
					corev1.ServicePort().WithName("infinispan").WithPort(11222),
				),
		)

	if err := ctx.Client().Apply(service); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan Service: %w", err))
	}
}

func ConfigMap(c *v1alpha1.Cache, ctx reconcile.Context) {
	config, err := configuration.Generate(&configuration.Spec{})
	if err != nil {
		ctx.Requeue(fmt.Errorf("unable to generate Infinispan configuration: %w", err))
		return
	}

	cm := corev1.
		ConfigMap(c.Name, c.Namespace).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithData(
			map[string]string{
				"infinispan.xml": config,
			},
		)

	if err := ctx.Client().Apply(cm); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan ConfigMap: %w", err))
	}
}

func DaemonSet(c *v1alpha1.Cache, ctx reconcile.Context) {
	ds := appsv1.
		DaemonSet(c.Name, c.Namespace).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DaemonSetSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(selectorLabels),
			).
			WithTemplate(corev1.PodTemplateSpec().
				WithName(containerName).
				WithLabels(selectorLabels).
				WithSpec(corev1.PodSpec().
					WithContainers(corev1.Container().
						WithName(containerName).
						WithImage("quay.io/infinispan/server:14.0").
						WithArgs("-c", "/config/infinispan.xml").
						WithPorts(
							corev1.ContainerPort().WithContainerPort(11222),
						).
						WithEnv(
							corev1.EnvVar().WithName("USER").WithValue("admin"),
							corev1.EnvVar().WithName("PASS").WithValue("password"),
						).
						WithLivenessProbe(
							httpProbe(5, 0, 10, 1, 80),
						).
						WithReadinessProbe(
							httpProbe(5, 0, 10, 1, 80),
						).
						WithStartupProbe(
							httpProbe(600, 1, 1, 1, 80),
						).
						WithVolumeMounts(
							corev1.VolumeMount().WithName("config").WithMountPath("/config").WithReadOnly(true),
						),
					).
					WithVolumes(corev1.Volume().
						WithName("config").
						WithConfigMap(
							corev1.ConfigMapVolumeSource().WithName(c.Name),
						),
					),
				),
			),
		)
	if err := ctx.Client().Apply(ds); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan DaemonSet: %w", err))
	}
}

func httpProbe(failureThreshold, initialDelay, period, successThreshold, timeout int32) *corev1.ProbeApplyConfiguration {
	return corev1.Probe().
		WithHTTPGet(
			corev1.HTTPGetAction().
				WithScheme(apicorev1.URISchemeHTTP).
				WithPath("rest/v2/cache-managers/default/health/status").
				WithPort(intstr.FromInt(11222)),
		).
		WithFailureThreshold(failureThreshold).
		WithInitialDelaySeconds(initialDelay).
		WithPeriodSeconds(period).
		WithSuccessThreshold(successThreshold).
		WithTimeoutSeconds(timeout)
}
