package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile/region"
	corev1 "k8s.io/api/core/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// TODO how to correctly handle side effects with dryRun?
// +kubebuilder:rbac:groups=core,namespace=gingersnap-operator-system,resources=pods;configmaps,verbs=create;delete;deletecollection;get;list;patch;update;watch

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create,versions=v1,name=sidecar.gingersnap-project.io,admissionReviewVersions=v1,sideEffects=None

const (
	ContainerName  = "gingersnap"
	ContainerImage = "registry.access.redhat.com/ubi9/ubi-minimal"
	VolumeName     = "gingersnap-regions"
	VolumeMount    = "/regions"
)

// ProxyInjector adds a cache side-car to pods with sidecar.gingersnap-project.io/inject: "true"
type ProxyInjector struct {
	Client  runtimeClient.Client
	decoder *admission.Decoder
}

// InjectDecoder injects the decoder.
func (injector *ProxyInjector) InjectDecoder(d *admission.Decoder) error {
	injector.decoder = d
	return nil
}

func (injector *ProxyInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	reqLogger := log.FromContext(ctx)

	pod := &corev1.Pod{}
	err := injector.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if !v1alpha1.CacheServiceLabelsExist(pod.Labels) {
		// TODO set higher level log
		reqLogger.Info("Cache labels don't exist, ignoring pod", "name", pod.Name, "generatename", pod.GenerateName, "namespace", pod.Namespace)
		return admission.Allowed("")
	}

	k8sClient := &client.Runtime{
		Client:    injector.Client,
		Ctx:       ctx,
		Namespace: req.Namespace,
		Scheme:    injector.Client.Scheme(),
	}

	cacheService := v1alpha1.CacheServiceFromLabels(pod.Labels)
	if err := k8sClient.WithNamespace(cacheService.Namespace).Load(cacheService.Name, &v1alpha1.Cache{}); err != nil {
		return admission.Errored(http.StatusBadRequest, fmt.Errorf("unable to load Cache %s: %w", cacheService, err))
	}

	configMap, err := injector.initRegionConfigMap(k8sClient, cacheService, pod.Namespace)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	addInlineProxyContainer(pod, configMap)
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// initRegionConfigMap populates a ConfigMap containing all CacheRegion definitions in the namespace of the Pod being admitted
func (injector *ProxyInjector) initRegionConfigMap(k8sClient client.Client, cacheService v1alpha1.CacheService, podNamespace string) (*corev1.ConfigMap, error) {
	configMap, err := region.FindExistingConfigMap(podNamespace, cacheService, k8sClient)
	if err != nil {
		return nil, fmt.Errorf("unable to search for existing ConfigMap: %v", err)
	}

	// ConfigMap doesn't exist yet, so we must create it
	if configMap != nil {
		return configMap, nil
	}
	return region.CreateConfigMap(podNamespace, cacheService, k8sClient)
}

func addInlineProxyContainer(pod *corev1.Pod, configMap *corev1.ConfigMap) {
	// TODO pass regions as args to proxy binary
	// If no region names passed, apply all configurations in --region-config-dir flag
	regions := v1alpha1.CacheRegionsFromAnnotations(pod.Annotations)

	pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
		Name:            ContainerName,
		Image:           ContainerImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command: []string{
			"sh",
			"-c",
			"while true; do ls -l /regions; echo \"\\n\"; sleep 10; done",
		},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      VolumeName,
			MountPath: VolumeMount,
		}},
	})

	// Only mount the requested CacheRegion configurations
	configMapItems := make([]corev1.KeyToPath, len(regions))
	for i, r := range regions {
		configMapItems[i] = corev1.KeyToPath{Key: r, Path: r}
	}

	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: VolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				Items: configMapItems,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMap.Name,
				},
			},
		},
	})
}
