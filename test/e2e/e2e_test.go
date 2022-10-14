//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"strings"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	k8s "github.com/gingersnap-project/operator/pkg/kubernetes"
	"github.com/gingersnap-project/operator/pkg/reconcile/sidecar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
		Expect(k8sClient.DeleteAllForeground(nil, &v1alpha1.CacheRegion{})).Should(Succeed())
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
			regionList := &v1alpha1.CacheRegionList{}
			Expect(k8sClient.List(nil, regionList)).Should(Succeed())
			return len(regionList.Items)
		}, Timeout, Interval).Should(Equal(0))
	})

	Context("Infinispan Deployment", func() {
		It("DaemonSet should be deployed successfully", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{
					Infinispan: &v1alpha1.InfinispanSpec{},
				},
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

	Context("Sidecar Injection", func() {
		It("region configmap should be removed from namespace when no pods exist", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())
			Eventually(func() error {
				return k8sClient.Load(cache.Name, cache)
			}, Timeout, Interval).Should(Succeed())

			objectMeta := meta("sidecar-injection")
			cache.CacheService().ApplyLabels(objectMeta.Labels)
			pod1 := &corev1.Pod{
				ObjectMeta: objectMeta,
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "app-container",
						Image:           "registry.access.redhat.com/ubi9/ubi-minimal",
						Args:            []string{"sleep", "1000000000"},
						ImagePullPolicy: corev1.PullIfNotPresent,
					}},
				},
			}
			pod2 := pod1.DeepCopy()
			pod2.Name += "-2"
			Expect(k8sClient.Create(pod1)).Should(Succeed())
			Expect(k8sClient.Create(pod2)).Should(Succeed())

			Expect(k8sClient.Load(pod1.Name, pod1)).Should(Succeed())
			configMapName := k8s.Volume(sidecar.VolumeName, &pod1.Spec).ConfigMap.Name

			configMap := &corev1.ConfigMap{}
			Eventually(func() []metav1.OwnerReference {
				Expect(k8sClient.Load(configMapName, configMap)).Should(Succeed())
				return configMap.OwnerReferences
			}, Timeout, Interval).Should(HaveLen(2))

			Expect(k8sClient.Delete(pod1.Name, pod1)).Should(Succeed())
			Expect(k8sClient.Delete(pod2.Name, pod2)).Should(Succeed())

			Eventually(func() bool {
				err := k8sClient.Load(configMapName, configMap)
				return errors.IsNotFound(err)
			}, Timeout, Interval).Should(BeTrue())
		})

		It("should inject proxy sidecar container when labels present", func() {

			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())

			objectMeta := meta("sidecar-injection")
			cache.CacheService().ApplyLabels(objectMeta.Labels)
			pod := &corev1.Pod{
				ObjectMeta: objectMeta,
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "app-container",
						Image:           "registry.access.redhat.com/ubi9/ubi-minimal",
						Args:            []string{"sleep", "1000000000"},
						ImagePullPolicy: corev1.PullIfNotPresent,
					}},
				},
			}
			Expect(k8sClient.Create(pod)).Should(Succeed())

			err := k8sClient.Load(pod.Name, pod)
			Expect(err).NotTo(HaveOccurred())

			configMapName := k8s.Volume(sidecar.VolumeName, &pod.Spec).ConfigMap.Name
			configMap := &corev1.ConfigMap{}
			Expect(k8sClient.Load(configMapName, configMap)).Should(Succeed())

			Expect(pod.Spec.Containers).Should(HaveLen(2))
			container := pod.Spec.Containers[1]
			Expect(container.Name).Should(Equal(sidecar.ContainerName))
			Expect(container.Image).Should(Equal(sidecar.ContainerImage))

			Expect(k8sClient.Delete(pod.Name, pod)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Load(configMapName, configMap)
				return errors.IsNotFound(err)
			}, Timeout, Interval).Should(BeTrue())
		})
	})

	Context("CacheRegion Update", func() {
		It("should propagate region changes to existing ConfigMaps", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())

			objectMeta := meta("propogate-region-changes")
			cache.CacheService().ApplyLabels(objectMeta.Labels)
			pod := &corev1.Pod{
				ObjectMeta: objectMeta,
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "app-container",
						Image:           "registry.access.redhat.com/ubi9/ubi-minimal",
						Args:            []string{"sleep", "1000000000"},
						ImagePullPolicy: corev1.PullIfNotPresent,
					}},
				},
			}
			Expect(k8sClient.Create(pod)).Should(Succeed())

			err := k8sClient.Load(pod.Name, pod)
			Expect(err).NotTo(HaveOccurred())

			configMapName := k8s.Volume(sidecar.VolumeName, &pod.Spec).ConfigMap.Name
			configMap := &corev1.ConfigMap{}
			Expect(k8sClient.Load(configMapName, configMap)).Should(Succeed())

			region := &v1alpha1.CacheRegion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "region-1",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheRegionSpec{
					Cache: cache.CacheService(),
				},
			}
			region2 := region.DeepCopy()
			region2.Name = "region-2"

			Expect(k8sClient.Create(region)).Should(Succeed())

			// Assert that CacheRegion is added to mounted ConfigMap
			Eventually(func() bool {
				Expect(k8sClient.Load(configMapName, configMap)).Should(Succeed())
				_, regionExists := configMap.BinaryData[region.Filename()]
				return regionExists
			}, Timeout, Interval).Should(BeTrue())

			Expect(k8sClient.Create(region2)).Should(Succeed())

			// Assert that subsequent CacheRegion is added to mounted ConfigMap
			Eventually(func() bool {
				Expect(k8sClient.Load(configMapName, configMap)).Should(Succeed())
				_, regionExists := configMap.BinaryData[region.Filename()]
				return regionExists
			}, Timeout, Interval).Should(BeTrue())
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
