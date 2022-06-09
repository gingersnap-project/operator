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

var _ = Describe("Cache Webhooks", func() {

	const timeout = time.Second * 30
	const interval = time.Second * 1

	key := types.NamespacedName{
		Name:      "cache-envtest",
		Namespace: "default",
	}

	AfterEach(func() {
		// Delete created resources
		By("Expecting to delete successfully")
		Eventually(func() error {
			f := &Cache{}
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
			f := &Cache{}
			return k8sClient.Get(ctx, key, f)
		}, timeout, interval).ShouldNot(Succeed())
	})

	It("should correctly set Cache defaults", func() {

		created := &Cache{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())

		Expect(k8sClient.Get(ctx, key, created)).Should(Succeed())
		spec := created.Spec
		// Ensure default values correctly set
		Expect(spec.Infinispan).ShouldNot(BeNil())
		Expect(spec.Redis).Should(BeNil())
	})

	It("should ensure that Cache cannot be created with both Infinispan and Redis specs", func() {
		rejected := &Cache{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: CacheSpec{
				Infinispan: &InfinispanSpec{},
				Redis:      &RedisSpec{},
			},
		}

		cause := statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec", "At most one of ['spec.infinispan', 'spec.redis'] must be configured"}
		expectInvalidErrStatus(k8sClient.Create(ctx, rejected), cause)
	})

	It("should ensure that Cache cannot be updated to contain both Infinispan and Redis specs", func() {

		created := &Cache{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: CacheSpec{
				Redis: &RedisSpec{},
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())

		// Ensure Spec is immutable on update
		updated := &Cache{}

		Expect(k8sClient.Get(ctx, key, updated)).Should(Succeed())
		updated.Spec.Infinispan = &InfinispanSpec{}

		cause := statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec", "At most one of ['spec.infinispan', 'spec.redis'] must be configured"}
		expectInvalidErrStatus(k8sClient.Update(ctx, updated), cause)
	})
})

type statusDetailCause struct {
	Type          metav1.CauseType
	field         string
	messageSubStr string
}

func expectInvalidErrStatus(err error, causes ...statusDetailCause) {
	Expect(err).ShouldNot(BeNil())
	var statusError *apierrors.StatusError
	Expect(errors.As(err, &statusError)).Should(BeTrue())

	errStatus := statusError.ErrStatus
	Expect(errStatus.Reason).Should(Equal(metav1.StatusReasonInvalid))

	Expect(errStatus.Details.Causes).Should(HaveLen(len(causes)))
	for i, c := range errStatus.Details.Causes {
		Expect(c.Type).Should(Equal(causes[i].Type))
		Expect(c.Field).Should(Equal(causes[i].field))
		Expect(c.Message).Should(ContainSubstring(causes[i].messageSubStr))
	}
}
