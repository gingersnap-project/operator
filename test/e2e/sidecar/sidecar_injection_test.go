//go:build e2e
// +build e2e

package sidecar

import (
	"context"
	"time"

	engytita "github.com/engytita/engytita-operator/pkg/reconcile/sidecar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("MutatingWebhook", func() {

	const timeout = time.Second * 30
	const interval = time.Second * 1

	ctx := context.TODO()
	meta := metav1.ObjectMeta{
		Name:      "sidecar-injection",
		Namespace: namespace,
		Labels: map[string]string{
			"app": "sidecar",
		},
	}

	AfterEach(func() {
		// Delete created resources
		By("Expecting to delete successfully")
		deleteOpts := []client.DeleteAllOfOption{
			client.InNamespace(meta.Namespace),
			client.MatchingLabels(meta.Labels),
		}
		Expect(k8sClient.DeleteAllOf(ctx, &corev1.Pod{}, deleteOpts...)).Should(Succeed())

		By("Expecting to delete finish")
		Eventually(func() int {
			podList := &corev1.PodList{}
			listOps := &client.ListOptions{Namespace: meta.Namespace, LabelSelector: labels.SelectorFromSet(meta.Labels)}
			Expect(k8sClient.List(ctx, podList, listOps)).Should(Succeed())
			return len(podList.Items)
		}, timeout, interval).Should(Equal(0))
	})

	Context("User created pod", func() {
		It("should inject proxy sidecar container when annotation true", func() {
			objectMeta := meta
			objectMeta.Annotations = map[string]string{engytita.AnnotationInject: "true"}
			pod := &corev1.Pod{
				ObjectMeta: objectMeta,
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "app-container",
						Image: "hello-world",
					}},
				},
			}
			Expect(k8sClient.Create(ctx, pod)).Should(Succeed())

			err := k8sClient.Get(ctx, client.ObjectKey{Namespace: meta.Namespace, Name: meta.Name}, pod)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod.Spec.Containers).Should(HaveLen(2))
			sidecar := pod.Spec.Containers[1]
			Expect(sidecar.Name).Should(Equal(engytita.ContainerName))
			Expect(sidecar.Image).Should(Equal(engytita.ContainerImage))
		})
	})
})
