package sidecar

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	k8s "github.com/engytita/engytita-operator/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=sidecar.engytita.org,admissionReviewVersions=v1,sideEffects=None

const (
	AnnotationInject = "sidecar.engytita.org/inject"
	ContainerName    = "engytita"
	ContainerImage   = "hello-world"
)

// ProxyInjector adds a cache side-car to pods with sidecar.engytita.org/inject: "true"
type ProxyInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *ProxyInjector) Handle(_ context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if err = injectProxyContainer(pod); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// InjectDecoder injects the decoder.
func (a *ProxyInjector) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}

func injectProxyContainer(pod *corev1.Pod) error {
	injectAnnotation, ok := pod.Annotations[AnnotationInject]
	if ok {
		inject, err := strconv.ParseBool(injectAnnotation)
		if err != nil {
			return err
		}

		if inject {
			container := k8s.GetContainer(ContainerName, &pod.Spec)
			// If the container already exists, we set the .image spec to the latest image
			if container != nil {
				container.Image = ContainerImage
			} else {
				pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
					Name:  ContainerName,
					Image: ContainerImage,
				})
			}
		}
	}
	return nil
}
