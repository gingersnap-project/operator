package mutation

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=sidecar.engytita.org,admissionReviewVersions=v1,sideEffects=None

// ProxyInjector adds a cache side-car to pods with sidecar.engytita.org/inject: "true"
type ProxyInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *ProxyInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	injectAnnotation, ok := pod.Annotations["sidecar.engytita.org/inject"]
	if ok {
		inject, err := strconv.ParseBool(injectAnnotation)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}

		if inject {
			addSideCarContainer(pod)
		}
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

func addSideCarContainer(pod *corev1.Pod) {
	container := corev1.Container{
		Name:  "proxy-sidecar",
		Image: "hello-world",
	}
	for i, c := range pod.Spec.Containers {
		if c.Name == container.Name {
			pod.Spec.Containers[i] = container
			return
		}
	}
	pod.Spec.Containers = append(pod.Spec.Containers, container)
}
