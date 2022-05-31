package sidecar

import (
	"testing"

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

	Context("User created pod", func() {
		It("should inject  proxy sidecar container when annotation true", func() {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{AnnotationInject: "true"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{existingContainer},
				},
			}
			Expect(injectProxyContainer(pod)).Should(Succeed())
			Expect(pod.Spec.Containers).Should(HaveLen(2))
			Expect(pod.Spec.Containers[1].Name).Should(Equal(ContainerName))
			Expect(pod.Spec.Containers[1].Image).Should(Equal(ContainerImage))
		})

		It("should update proxy sidecar container image when sidecar already present", func() {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{AnnotationInject: "true"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						existingContainer,
						{
							Name:  ContainerName,
							Image: "some-out-dated-image",
						},
					},
				},
			}
			Expect(injectProxyContainer(pod)).Should(Succeed())
			Expect(pod.Spec.Containers).Should(HaveLen(2))
			Expect(pod.Spec.Containers[1].Name).Should(Equal(ContainerName))
			Expect(pod.Spec.Containers[1].Image).Should(Equal(ContainerImage))
		})

		It("should do nothing when annotation false", func() {
			pod := &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{existingContainer},
				},
			}
			Expect(injectProxyContainer(pod)).Should(Succeed())
			Expect(pod.Spec.Containers).Should(HaveLen(1))
			Expect(pod.Spec.Containers[0].Name).Should(Equal(existingContainer.Name))
		})

		It("should do nothing when annotation not present", func() {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{AnnotationInject: "false"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{existingContainer},
				},
			}
			Expect(injectProxyContainer(pod)).Should(Succeed())
			Expect(pod.Spec.Containers).Should(HaveLen(1))
			Expect(pod.Spec.Containers[0].Name).Should(Equal(existingContainer.Name))
		})
	})
})
