package redis

import (
	"fmt"
	"strconv"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache/context"
	apicorev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const (
	containerName = "redis"
)

var (
	selectorLabels = map[string]string{
		"app": "redis",
	}
)

func Service(c *v1alpha1.Cache, ctx *context.Context) {
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
					corev1.ServicePort().WithName("redis").WithPort(6379),
				),
		)

	if err := ctx.Client().Apply(service); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Redis Service: %w", err))
	}
}

func ConfigurationSecret(c *v1alpha1.Cache, ctx *context.Context) {
	secretName := c.ConfigurationSecret()

	// Initialize the ctx ServiceBinding so that we can use the values when creating the DaemonSet
	sb := &context.ServiceBinding{
		Port: 6379,
		Host: c.Name,
	}
	ctx.ServiceBinding = sb

	secret := corev1.Secret(secretName, c.Namespace).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithStringData(
			map[string]string{
				"type":     "redis",
				"provider": "gingersnap",
				"host":     sb.Host,
				"port":     strconv.Itoa(sb.Port),
			},
		).
		WithType("servicebinding.io/redis")

	if err := ctx.Client().Apply(secret); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Redis configuration secret: %w", err))
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
	sb := ctx.ServiceBinding
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
						WithImage("redis:7.0.0").
						WithPorts(
							corev1.ContainerPort().WithContainerPort(int32(sb.Port)),
						),
					),
				),
			),
		)
	if err := ctx.Client().Apply(ds); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Redis DaemonSet: %w", err))
	}
}
