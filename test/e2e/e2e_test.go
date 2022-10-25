//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"strings"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const Timeout = time.Second * 60 * 4
const Interval = time.Second * 1

var _ = Describe("E2E", func() {

	AfterEach(func() {
		meta := meta("")
		test := CurrentGinkgoTestDescription()
		if test.Failed {
			dir := fmt.Sprintf("%s/%s", OutputDir, strings.ReplaceAll(test.TestText, " ", "_"))
			k8sClient.WriteAllResourcesToFile(dir)
		}

		// Delete created test resources
		By("Expecting to delete successfully")
		// Delete all CRs in the foreground to ensure that any dependent resources are deleted before the resource
		// This simplifies the logic below as it's not necessary to check that all subordinate resource types have been
		// removed from the namespace
		Expect(k8sClient.DeleteAllForeground(nil, &v1alpha1.Cache{})).Should(Succeed())
		Expect(k8sClient.DeleteAllForeground(nil, &v1alpha1.LazyCacheRule{})).Should(Succeed())
		Expect(k8sClient.DeleteAllOf(meta.Labels, &corev1.Pod{})).Should(Succeed())

		By("Expecting delete to finish")
		Eventually(func() int {
			podList := &corev1.PodList{}
			Expect(k8sClient.List(meta.Labels, podList)).Should(Succeed())
			return len(podList.Items)
		}, Timeout, Interval).Should(Equal(0))

		Eventually(func() int {
			cacheList := &v1alpha1.CacheList{}
			Expect(k8sClient.List(nil, cacheList)).Should(Succeed())
			return len(cacheList.Items)
		}, Timeout, Interval).Should(Equal(0))

		Eventually(func() int {
			ruleList := &v1alpha1.LazyCacheRuleList{}
			Expect(k8sClient.List(nil, ruleList)).Should(Succeed())
			return len(ruleList.Items)
		}, Timeout, Interval).Should(Equal(0))
	})

	Context("Cache Deployment", func() {
		It("DaemonSet should be deployed successfully", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())

			secret := &corev1.Secret{}
			Eventually(func() error {
				return k8sClient.Load(cache.ConfigurationSecret(), secret)
			}, Timeout, Interval).Should(Succeed())

			Expect(secret.Data).To(HaveKeyWithValue("type", []byte("infinispan")))
			Expect(secret.Data).To(HaveKeyWithValue("provider", []byte("gingersnap")))
			Expect(secret.Data).To(HaveKeyWithValue("host", []byte(cache.Name)))
			Expect(secret.Data).To(HaveKeyWithValue("username", []byte("admin")))
			Expect(secret.Data).To(HaveKeyWithValue("port", []byte("11222")))
			Expect(secret.Data).To(HaveKey("password"))
			Expect(secret.Type).Should(Equal(corev1.SecretType("servicebinding.io/infinispan")))

			Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
			Expect(cache.Status.ServiceBinding.Name).Should(Equal(cache.ConfigurationSecret()))

			daemonSet := &appsv1.DaemonSet{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, daemonSet)
			}, Timeout, Interval).Should(Succeed())

			Eventually(func() bool {
				Expect(k8sClient.Load(cache.Name, daemonSet)).Should(Succeed())
				return daemonSet.Status.CurrentNumberScheduled > 0 && daemonSet.Status.NumberUnavailable == 0
			}, Timeout, Interval).Should(BeTrue())
		})
	})

	Context("LazyCacheRule", func() {
		It("Cache ConfigMap should be created with rule", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())

			cacheRule := &v1alpha1.LazyCacheRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "lazy-cache-rule",
					Namespace: Namespace,
				},
				Spec: v1alpha1.LazyCacheRuleSpec{
					Cache: v1alpha1.CacheService{
						Name:      cache.Name,
						Namespace: cache.Namespace,
					},
				},
			}
			Expect(k8sClient.Create(cacheRule)).Should(Succeed())

			cm := &corev1.ConfigMap{}
			cmName := cache.CacheService().LazyCacheConfigMap()
			Eventually(func() error {
				return k8sClient.Load(cmName, cm)
			}, Timeout, Interval).Should(Succeed())

			Expect(cm.BinaryData).Should(HaveLen(1))
			Expect(cm.BinaryData).To(HaveKey(cacheRule.Filename()))

			Expect(k8sClient.Delete(cacheRule.Name, cacheRule)).Should(Succeed())
			Eventually(func() map[string][]byte {
				_ = k8sClient.Load(cmName, cm)
				return cm.BinaryData
			}, Timeout, Interval).Should(Not(HaveKey(cacheRule.Filename())))
		})
	})
})

func meta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: Namespace,
		Labels: map[string]string{
			"app": "e2e-test",
		},
	}
}
