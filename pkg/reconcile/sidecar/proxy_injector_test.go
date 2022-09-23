package sidecar

import (
	"testing"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

func TestSidecar(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Sidecar Unit",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = Describe("ProxyInjector", func() {

	existingContainer := corev1.Container{
		Name:  "existing",
		Image: "some-image",
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "region-configmap",
		},
	}

	Context("User created pod", func() {
		It("should inject  proxy sidecar container when annotation true", func() {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{v1alpha1.AnnotationRegions: "Region1"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{existingContainer},
				},
			}

			addInlineProxyContainer(pod, configMap)
			Expect(pod.Spec.Containers).Should(HaveLen(2))
			Expect(pod.Spec.Containers[1].Name).Should(Equal(ContainerName))
			Expect(pod.Spec.Containers[1].Image).Should(Equal(ContainerImage))
			Expect(pod.Spec.Containers[1].VolumeMounts).Should(HaveLen(1))
			Expect(pod.Spec.Containers[1].VolumeMounts[0].Name).Should(Equal(VolumeName))
			Expect(pod.Spec.Containers[1].VolumeMounts[0].MountPath).Should(Equal(VolumeMount))
			Expect(pod.Spec.Volumes).Should(HaveLen(1))
			Expect(pod.Spec.Volumes[0].Name).Should(Equal(VolumeName))
			Expect(pod.Spec.Volumes[0].ConfigMap.Name).Should(Equal(configMap.Name))
			Expect(pod.Spec.Volumes[0].ConfigMap.Items).Should(Equal([]corev1.KeyToPath{{Key: "Region1", Path: "Region1"}}))
		})
	})
})
