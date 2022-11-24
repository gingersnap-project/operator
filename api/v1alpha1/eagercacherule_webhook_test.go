package v1alpha1

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("EagerCacheRule Webhooks", func() {

	const timeout = time.Second * 30
	const interval = time.Second * 1

	key := types.NamespacedName{
		Name:      "eagercacherule-envtest",
		Namespace: "default",
	}

	AfterEach(func() {
		// Delete created resources
		By("Expecting to delete successfully")
		Eventually(func() error {
			f := &EagerCacheRule{}
			if err := k8sClient.Get(ctx, key, f); err != nil {
				var statusError *apierrors.StatusError
				if !errors.As(err, &statusError) {
					return err
				}
				// If the resource does not exist, do nothing
				if statusError.ErrStatus.Code == 404 {
					return nil
				}
			}
			return k8sClient.Delete(ctx, f)
		}, timeout, interval).Should(Succeed())

		By("Expecting to delete finish")
		Eventually(func() error {
			f := &EagerCacheRule{}
			return k8sClient.Get(ctx, key, f)
		}, timeout, interval).ShouldNot(Succeed())
	})

	It("should correctly set EagerCacheRule defaults", func() {

		created := &EagerCacheRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: EagerCacheRuleSpec{
				CacheRef: &NamespacedRef{
					Name:      "cache1",
					Namespace: "cache2",
				},
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())

		Expect(k8sClient.Get(ctx, key, created)).Should(Succeed())
		Expect(created.Labels).To(HaveKeyWithValue("gingersnap-project.io/cache", created.Spec.CacheRef.Name))
		Expect(created.Labels).To(HaveKeyWithValue("gingersnap-project.io/cache-namespace", created.Spec.CacheRef.Namespace))
	})
})