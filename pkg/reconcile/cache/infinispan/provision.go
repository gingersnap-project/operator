package infinispan

import (
	"fmt"
	"strconv"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/infinispan/configuration"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache/context"
	"github.com/gingersnap-project/operator/pkg/reconcile/meta"
	"github.com/gingersnap-project/operator/pkg/security/passwords"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

const (
	containerName = "infinispan"
)

var (
	labels = meta.GingersnapLabels("infinispan", "cache")
)

func Service(c *v1alpha1.Cache, ctx *context.Context) {
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
					corev1.ServicePort().WithName("infinispan").WithPort(11222),
				),
		)

	if err := ctx.Client().Apply(service); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply Infinispan Service: %w", err))
	}
}

func ConfigMap(c *v1alpha1.Cache, ctx *context.Context) {
	config, err := configuration.Generate(&configuration.Spec{})
	if err != nil {
		ctx.Requeue(fmt.Errorf("unable to generate Infinispan configuration: %w", err))
		return
	}

	cm := corev1.
		ConfigMap(c.Name, c.Namespace).
		WithLabels(labels).
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

func ConfigurationSecret(c *v1alpha1.Cache, ctx *context.Context) {
	secretName := c.ConfigurationSecret()
	existingSecret := &apicorev1.Secret{}
	var password string
	if err := ctx.Client().Load(secretName, existingSecret); err != nil {
		if !errors.IsNotFound(err) {
			ctx.Requeue(fmt.Errorf("unable to retrieve existing configuration secret: %w", err))
			return
		}

		password, err = passwords.Generate(16)
		if err != nil {
			ctx.Requeue(fmt.Errorf("unable to generate password: %w", err))
			return
		}
	} else {
		// Configuration secret already exists, so make sure we apply the existing password
		password = string(existingSecret.Data["password"])
	}

	// Initialize the ctx ServiceBinding so that we can use the values when creating the DaemonSet
	sb := &context.ServiceBinding{
		Username: "admin",
		Password: password,
		Port:     11222,
		Host:     c.Name,
	}
	ctx.ServiceBinding = sb

	secret := corev1.Secret(secretName, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(
			ctx.Client().OwnerReference(),
		).
		WithStringData(
			map[string]string{
				"type":     "infinispan",
				"provider": "gingersnap",
				"host":     sb.Host,
				"port":     strconv.Itoa(sb.Port),
				"username": sb.Username,
				"password": sb.Password,
			},
		).
		WithType("servicebinding.io/infinispan")

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
	sb := ctx.ServiceBinding
	ds := appsv1.
		DaemonSet(c.Name, c.Namespace).
		WithLabels(labels).
		WithOwnerReferences(ctx.Client().OwnerReference()).
		WithSpec(appsv1.DaemonSetSpec().
			WithSelector(
				metav1.LabelSelector().WithMatchLabels(labels),
			).
			WithTemplate(corev1.PodTemplateSpec().
				WithName(containerName).
				WithLabels(labels).
				WithSpec(corev1.PodSpec().
					WithContainers(corev1.Container().
						WithName(containerName).
						WithImage("quay.io/infinispan/server:14.0").
						WithArgs("-c", "/config/infinispan.xml").
						WithPorts(
							corev1.ContainerPort().WithContainerPort(int32(sb.Port)),
						).
						WithEnv(
							corev1.EnvVar().WithName("USER").WithValue(sb.Username),
							corev1.EnvVar().WithName("PASS").WithValue(sb.Password),
						).
						WithLivenessProbe(
							httpProbe(5, 0, 10, 1, 80, sb.Port),
						).
						WithReadinessProbe(
							httpProbe(5, 0, 10, 1, 80, sb.Port),
						).
						WithStartupProbe(
							httpProbe(600, 1, 1, 1, 80, sb.Port),
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

func httpProbe(failureThreshold, initialDelay, period, successThreshold, timeout int32, port int) *corev1.ProbeApplyConfiguration {
	return corev1.Probe().
		WithHTTPGet(
			corev1.HTTPGetAction().
				WithScheme(apicorev1.URISchemeHTTP).
				WithPath("rest/v2/cache-managers/default/health/status").
				WithPort(intstr.FromInt(port)),
		).
		WithFailureThreshold(failureThreshold).
		WithInitialDelaySeconds(initialDelay).
		WithPeriodSeconds(period).
		WithSuccessThreshold(successThreshold).
		WithTimeoutSeconds(timeout)
}
