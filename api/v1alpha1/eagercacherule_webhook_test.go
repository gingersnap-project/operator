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
				CacheRef: &NamespacedObjectReference{
					Name:      "cache1",
					Namespace: "cache2",
				},
				TableName: "SomeTable",
				Key: &EagerCacheKey{
					KeyColumns: []string{"col1"},
				},
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())

		Expect(k8sClient.Get(ctx, key, created)).Should(Succeed())
		Expect(created.Labels).To(HaveKeyWithValue("gingersnap-project.io/cache", created.Spec.CacheRef.Name))
		Expect(created.Labels).To(HaveKeyWithValue("gingersnap-project.io/cache-namespace", created.Spec.CacheRef.Namespace))
		Expect(created.Spec.Key.Format).Should(Equal(KeyFormat_TEXT))
		Expect(created.Spec.Key.KeySeparator).Should(Equal("|"))
	})

	It("Should reject rule if required fields are missing", func() {
		invalid := &EagerCacheRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: EagerCacheRuleSpec{},
		}

		ExpectInvalidErrStatus(k8sClient.Create(ctx, invalid),
			statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec.cacheRef", "'cacheRef' field must be defined"},
			statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec.tableName", "'tableName' field must not be empty"},
			statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec.key", "'key' field must be defined"},
		)

		invalid.Spec.CacheRef = &NamespacedObjectReference{}
		invalid.Spec.TableName = "Table"
		invalid.Spec.Key = &EagerCacheKey{}
		ExpectInvalidErrStatus(k8sClient.Create(ctx, invalid),
			statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec.cacheRef.name", "'name' field must not be empty"},
			statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec.cacheRef.namespace", "'namespace' field must not be empty"},
			statusDetailCause{metav1.CauseTypeFieldValueRequired, "spec.key.keyColumns", "'keyColumns' field must not be empty"},
		)
	})

	It("Should return error if any spec value is updated", func() {

		created := &EagerCacheRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: EagerCacheRuleSpec{
				CacheRef: &NamespacedObjectReference{
					Name:      "cache1",
					Namespace: "cache2",
				},
				TableName: "SomeTable",
				Key: &EagerCacheKey{
					KeyColumns: []string{"col1"},
				},
				Value: &Value{
					ValueColumns: []string{"col1", "col2"},
				},
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())

		// Ensure that non-spec fields can still be updated
		created.ObjectMeta.Labels = map[string]string{"example": "label"}
		Expect(k8sClient.Update(ctx, created)).Should(Succeed())

		// Ensure Spec is immutable on update
		updated := &EagerCacheRule{}
		Expect(k8sClient.Get(ctx, key, updated)).Should(Succeed())

		cause := statusDetailCause{"FieldValueForbidden", "spec", "The EagerCacheRule spec is immutable and cannot be updated after initial creation"}
		updated.Spec.TableName = "New Table"
		ExpectInvalidErrStatus(k8sClient.Update(ctx, updated), cause)

		Expect(k8sClient.Get(ctx, key, updated)).Should(Succeed())
		updated.Spec.CacheRef = nil
		ExpectInvalidErrStatus(k8sClient.Update(ctx, updated), cause)

		Expect(k8sClient.Get(ctx, key, updated)).Should(Succeed())
		updated.Spec.Key = nil
		ExpectInvalidErrStatus(k8sClient.Update(ctx, updated), cause)

		Expect(k8sClient.Get(ctx, key, updated)).Should(Succeed())
		updated.Spec.Value = nil
		ExpectInvalidErrStatus(k8sClient.Update(ctx, updated), cause)
	})
})
