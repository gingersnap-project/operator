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

	It("should correctly set Local Cache defaults", func() {

		created := &Cache{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())
		Expect(k8sClient.Get(ctx, key, created)).Should(Succeed())
		Expect(created.Spec.Deployment.Type).Should(Equal(CacheDeploymentType_LOCAL))
		Expect(created.Spec.Deployment.Replicas).Should(Equal(int32(0)))
	})

	It("should correctly set Cluster Cache defaults", func() {

		created := &Cache{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: CacheSpec{
				Deployment: &CacheDeploymentSpec{
					Type: CacheDeploymentType_CLUSTER,
				},
			},
		}

		Expect(k8sClient.Create(ctx, created)).Should(Succeed())
		Expect(k8sClient.Get(ctx, key, created)).Should(Succeed())
		Expect(created.Spec.Deployment.Type).Should(Equal(CacheDeploymentType_CLUSTER))
		Expect(created.Spec.Deployment.Replicas).Should(Equal(int32(1)))
	})

	It("should reject invalid resource quantities", func() {

		valid := &Cache{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: CacheSpec{
				Deployment: &CacheDeploymentSpec{
					Resources: &Resources{
						Requests: &ResourceQuantity{
							Cpu:    "0.1",
							Memory: "512Mi",
						},
						Limits: &ResourceQuantity{
							Cpu:    "1",
							Memory: "512Mi",
						},
					},
				},
				DbSyncer: &DBSyncerDeploymentSpec{
					Resources: &Resources{
						Requests: &ResourceQuantity{
							Cpu:    "0.1",
							Memory: "512Mi",
						},
						Limits: &ResourceQuantity{
							Cpu:    "1",
							Memory: "512Mi",
						},
					},
				},
			},
		}

		Expect(k8sClient.Create(ctx, valid)).Should(Succeed())

		invalid := valid.DeepCopy()
		invalid.Spec.Deployment.Resources.Requests.Cpu = "regex fail"
		invalid.Spec.Deployment.Resources.Requests.Memory = "512mi"
		invalid.Spec.Deployment.Resources.Limits.Cpu = "regex fail"
		invalid.Spec.Deployment.Resources.Limits.Memory = "1a"
		invalid.Spec.DbSyncer.Resources.Requests.Cpu = "regex fail"
		invalid.Spec.DbSyncer.Resources.Requests.Memory = "512mi"
		invalid.Spec.DbSyncer.Resources.Limits.Cpu = "regex fail"
		invalid.Spec.DbSyncer.Resources.Limits.Memory = "1a"

		expectInvalidErrStatus(k8sClient.Create(ctx, invalid),
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.deployment.resources.requests.cpu", "quantities must match the regular expression"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.deployment.resources.requests.memory", "unable to parse quantity's suffix"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.deployment.resources.limits.cpu", "quantities must match the regular expression"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.deployment.resources.limits.memory", "quantities must match the regular expression"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.dbSyncer.resources.requests.cpu", "quantities must match the regular expression"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.dbSyncer.resources.requests.memory", "unable to parse quantity's suffix"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.dbSyncer.resources.limits.cpu", "quantities must match the regular expression"},
			statusDetailCause{metav1.CauseTypeFieldValueInvalid, "spec.dbSyncer.resources.limits.memory", "quantities must match the regular expression"},
		)
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
