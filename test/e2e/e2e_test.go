//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"strings"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	bindingv1 "github.com/gingersnap-project/operator/pkg/apis/binding/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
		Expect(k8sClient.DeleteAllForeground(nil, &v1alpha1.LazyCacheRule{})).Should(Succeed())
		Expect(k8sClient.DeleteAllForeground(nil, &v1alpha1.EagerCacheRule{})).Should(Succeed())
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

	Context("Local Cache Deployment", func() {
		It("DaemonSet should be deployed successfully", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{
					DataSource: &v1alpha1.DataSourceSpec{
						DbType: v1alpha1.DBType_MYSQL_8.Enum(),
						SecretRef: &v1alpha1.LocalObjectReference{
							Name: MysqlConnectionSecret.Name,
						},
					},
					Deployment: &v1alpha1.CacheDeploymentSpec{
						Resources: &v1alpha1.Resources{
							Requests: &v1alpha1.ResourceQuantity{
								Cpu:    "500m",
								Memory: "512Mi",
							},
							Limits: &v1alpha1.ResourceQuantity{
								Cpu:    "1",
								Memory: "1Gi",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
				return cache.Condition(v1alpha1.CacheConditionReady).Status == metav1.ConditionTrue
			}, Timeout, Interval).Should(BeTrue())

			// Ensure Cache ServiceBinding created correctly
			expectSBSecret(
				cache.CacheService().UserServiceBindingSecret(),
				cache.CacheService().SvcName(),
				"8080",
			)

			expectServiceBinding(
				cache.CacheService().DataSourceServiceBinding(),
				"mysql",
				cache.Spec.DataSource.SecretRef.Name,
				"DaemonSet",
				cache.Name,
			)

			// Ensure db-syncer Cache ServiceBinding secret created correctly
			expectSBSecret(
				cache.CacheService().DBSyncerCacheServiceBindingSecret(),
				cache.CacheService().SvcName(),
				"11222",
			)

			Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
			Expect(cache.Status.ServiceBinding.Name).Should(Equal(cache.CacheService().UserServiceBindingSecret()))

			sa := &corev1.ServiceAccount{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, sa)
			}, Timeout, Interval).Should(Succeed())

			role := &rbacv1.Role{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, role)
			}, Timeout, Interval).Should(Succeed())

			Expect(role.Rules[0].Resources).Should(ContainElement("configmaps"))
			Expect(role.Rules[0].Verbs).Should(ContainElements("get", "list", "watch"))

			roleBinding := &rbacv1.RoleBinding{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, roleBinding)
			}, Timeout, Interval).Should(Succeed())
			Expect(roleBinding.RoleRef.Name).Should(Equal(role.Name))
			Expect(roleBinding.Subjects[0].Name).Should(Equal(sa.Name))

			daemonSet := &appsv1.DaemonSet{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, daemonSet)
			}, Timeout, Interval).Should(Succeed())

			Eventually(func() bool {
				Expect(k8sClient.Load(cache.Name, daemonSet)).Should(Succeed())
				return daemonSet.Status.NumberReady > 0 && daemonSet.Status.NumberUnavailable == 0
			}, Timeout, Interval).Should(BeTrue())

			res := daemonSet.Spec.Template.Spec.Containers[0].Resources
			Expect(res.Requests.Cpu().String()).Should(Equal("500m"))
			Expect(res.Requests.Memory().String()).Should(Equal("512Mi"))
			Expect(res.Limits.Cpu().String()).Should(Equal("1"))
			Expect(res.Limits.Memory().String()).Should(Equal("1Gi"))
		})
	})

	Context("Cluster Cache Deployment", func() {
		It("Deployment should be deployed successfully", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{
					DataSource: &v1alpha1.DataSourceSpec{
						DbType: v1alpha1.DBType_MYSQL_8.Enum(),
						SecretRef: &v1alpha1.LocalObjectReference{
							Name: MysqlConnectionSecret.Name,
						},
					},
					Deployment: &v1alpha1.CacheDeploymentSpec{
						Type:     v1alpha1.CacheDeploymentType_CLUSTER,
						Replicas: 2,
						Resources: &v1alpha1.Resources{
							Requests: &v1alpha1.ResourceQuantity{
								Cpu:    "100m",
								Memory: "256Mi",
							},
							Limits: &v1alpha1.ResourceQuantity{
								Cpu:    "200m",
								Memory: "512Mi",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
				return cache.Condition(v1alpha1.CacheConditionReady).Status == metav1.ConditionTrue
			}, Timeout, Interval).Should(BeTrue())

			// Ensure Cache ServiceBinding created correctly
			expectSBSecret(
				cache.CacheService().UserServiceBindingSecret(),
				cache.CacheService().SvcName(),
				"8080",
			)

			expectServiceBinding(
				cache.CacheService().DataSourceServiceBinding(),
				"mysql",
				cache.Spec.DataSource.SecretRef.Name,
				"Deployment",
				cache.Name,
			)

			// Ensure db-syncer Cache ServiceBinding secret created correctly
			expectSBSecret(
				cache.CacheService().DBSyncerCacheServiceBindingSecret(),
				cache.CacheService().SvcName(),
				"11222",
			)

			Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
			Expect(cache.Status.ServiceBinding.Name).Should(Equal(cache.CacheService().UserServiceBindingSecret()))

			sa := &corev1.ServiceAccount{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, sa)
			}, Timeout, Interval).Should(Succeed())

			role := &rbacv1.Role{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, role)
			}, Timeout, Interval).Should(Succeed())

			Expect(role.Rules[0].Resources).Should(ContainElement("configmaps"))
			Expect(role.Rules[0].Verbs).Should(ContainElements("get", "list", "watch"))

			roleBinding := &rbacv1.RoleBinding{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, roleBinding)
			}, Timeout, Interval).Should(Succeed())
			Expect(roleBinding.RoleRef.Name).Should(Equal(role.Name))
			Expect(roleBinding.Subjects[0].Name).Should(Equal(sa.Name))

			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Load(cache.Name, deployment)
			}, Timeout, Interval).Should(Succeed())

			Eventually(func() int32 {
				Expect(k8sClient.Load(cache.Name, deployment)).Should(Succeed())
				return deployment.Status.AvailableReplicas
			}, Timeout, Interval).Should(Equal(int32(2)))

			res := deployment.Spec.Template.Spec.Containers[0].Resources
			Expect(res.Requests.Cpu().String()).Should(Equal("100m"))
			Expect(res.Requests.Memory().String()).Should(Equal("256Mi"))
			Expect(res.Limits.Cpu().String()).Should(Equal("200m"))
			Expect(res.Limits.Memory().String()).Should(Equal("512Mi"))
		})
	})

	Context("LazyCacheRule", func() {
		It("Cache ConfigMap should be created with rule", func() {
			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{
					DataSource: &v1alpha1.DataSourceSpec{
						DbType: v1alpha1.DBType_MYSQL_8.Enum(),
						SecretRef: &v1alpha1.LocalObjectReference{
							Name: MysqlConnectionSecret.Name,
						},
					},
				},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())

			cacheRule := &v1alpha1.LazyCacheRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "lazy-cache-rule",
					Namespace: Namespace,
				},
				Spec: v1alpha1.LazyCacheRuleSpec{
					CacheRef: &v1alpha1.NamespacedObjectReference{
						Name:      cache.Name,
						Namespace: cache.Namespace,
					},
					Query: "TODO replace with actual DB query",
				},
			}
			Expect(k8sClient.Create(cacheRule)).Should(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
				return cache.Condition(v1alpha1.CacheConditionReady).Status == metav1.ConditionTrue
			}, Timeout, Interval).Should(BeTrue())

			cm := &corev1.ConfigMap{}
			cmName := cacheRule.ConfigMap()
			Eventually(func() error {
				return k8sClient.Load(cmName, cm)
			}, Timeout, Interval).Should(Succeed())

			Expect(cm.Data).Should(HaveLen(1))
			Expect(cm.Data).To(HaveKey(cacheRule.GetName()))

			Expect(k8sClient.Delete(cacheRule.Name, cacheRule)).Should(Succeed())
			Eventually(func() map[string]string {
				_ = k8sClient.Load(cmName, cm)
				return cm.Data
			}, Timeout, Interval).Should(Not(HaveKey(cacheRule.GetName())))
		})
	})

	Context("EagerCacheRule", func() {
		// TODO add integration test for DataSource using ServiceProviderRef
		It("ConfigMap should be created with rule and db-syncer deployed", func() {

			cache := &v1alpha1.Cache{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cache",
					Namespace: Namespace,
				},
				Spec: v1alpha1.CacheSpec{
					DataSource: &v1alpha1.DataSourceSpec{
						DbType: v1alpha1.DBType_MYSQL_8.Enum(),
						SecretRef: &v1alpha1.LocalObjectReference{
							Name: MysqlConnectionSecret.Name,
						},
					},
					DbSyncer: &v1alpha1.DBSyncerDeploymentSpec{
						Resources: &v1alpha1.Resources{
							Requests: &v1alpha1.ResourceQuantity{
								Cpu:    "500m",
								Memory: "512Mi",
							},
							Limits: &v1alpha1.ResourceQuantity{
								Cpu:    "1",
								Memory: "1Gi",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(cache)).Should(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Load(cache.Name, cache)).Should(Succeed())
				return cache.Condition(v1alpha1.CacheConditionReady).Status == metav1.ConditionTrue
			}, Timeout, Interval).Should(BeTrue())

			cacheService := cache.CacheService()

			cacheRule := &v1alpha1.EagerCacheRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "eager-cache-rule",
					Namespace: Namespace,
				},
				Spec: v1alpha1.EagerCacheRuleSpec{
					CacheRef: &v1alpha1.NamespacedObjectReference{
						Name:      cache.Name,
						Namespace: cache.Namespace,
					},
					Key: &v1alpha1.EagerCacheKey{
						KeyColumns: []string{"id"},
					},
					TableName: "debezium.customer",
				},
			}
			Expect(k8sClient.Create(cacheRule)).Should(Succeed())

			cm := &corev1.ConfigMap{}
			cmName := cacheRule.ConfigMap()
			Eventually(func() error {
				return k8sClient.Load(cmName, cm)
			}, Timeout, Interval).Should(Succeed())

			Expect(cm.Data).Should(HaveLen(1))
			Expect(cm.Data).To(HaveKey(cacheRule.GetName()))

			// Ensure Cache ServiceBinding created correctly
			expectSBSecret(
				cache.CacheService().UserServiceBindingSecret(),
				cache.CacheService().SvcName(),
				"8080",
			)

			expectServiceBinding(
				cache.CacheService().DataSourceServiceBinding(),
				"mysql",
				cache.Spec.DataSource.SecretRef.Name,
				"DaemonSet",
				cache.Name,
			)

			// Ensure db-syncer Cache ServiceBinding secret created correctly
			expectSBSecret(
				cache.CacheService().DBSyncerCacheServiceBindingSecret(),
				cache.CacheService().SvcName(),
				"11222",
			)

			expectServiceBinding(
				cacheService.DBSyncerCacheServiceBinding(),
				"gingersnap",
				cacheService.DBSyncerCacheServiceBindingSecret(),
				"Deployment",
				cacheService.DBSyncerName(),
			)

			dbSyncer := &appsv1.Deployment{}
			Eventually(func() bool {
				if err := k8sClient.Load(cacheService.DBSyncerName(), dbSyncer); err != nil {
					return false
				}
				return dbSyncer.Status.ReadyReplicas == 1
			}, Timeout, Interval).Should(BeTrue())

			res := dbSyncer.Spec.Template.Spec.Containers[0].Resources
			Expect(res.Requests.Cpu().String()).Should(Equal("500m"))
			Expect(res.Requests.Memory().String()).Should(Equal("512Mi"))
			Expect(res.Limits.Cpu().String()).Should(Equal("1"))
			Expect(res.Limits.Memory().String()).Should(Equal("1Gi"))

			// Ensure all resources are cleaned up on rule deletion
			Expect(k8sClient.Delete(cacheRule.Name, cacheRule)).Should(Succeed())
			Eventually(func() map[string]string {
				_ = k8sClient.Load(cmName, cm)
				return cm.Data
			}, Timeout, Interval).Should(Not(HaveKey(cacheRule.GetName())))

			Eventually(func() bool {
				return errors.IsNotFound(k8sClient.Load(cacheRule.CacheService().DBSyncerName(), dbSyncer))
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

func expectSBSecret(name, svc, port string) {
	secret := &corev1.Secret{}
	Eventually(func() error {
		return k8sClient.Load(name, secret)
	}, Timeout, Interval).Should(Succeed())

	Expect(secret.Data).To(HaveKeyWithValue("type", []byte("gingersnap")))
	Expect(secret.Data).To(HaveKeyWithValue("provider", []byte("gingersnap")))
	Expect(secret.Data).To(HaveKeyWithValue("host", []byte(svc)))
	Expect(secret.Data).To(HaveKeyWithValue("port", []byte(port)))
	Expect(secret.Type).Should(Equal(corev1.SecretType("servicebinding.io/gingersnap")))
}

func expectServiceBinding(name, bindingType, secret, workloadKind, workloadName string) {
	sb := &bindingv1.ServiceBinding{}
	Eventually(func() error {
		return k8sClient.Load(name, sb)
	}, Timeout, Interval).Should(Succeed())
	Expect(sb.Spec.Type).Should(Equal(bindingType))
	Expect(sb.Spec.Service.APIVersion).Should(Equal("v1"))
	Expect(sb.Spec.Service.Kind).Should(Equal("Secret"))
	Expect(sb.Spec.Service.Name).Should(Equal(secret))
	Expect(sb.Spec.Workload.APIVersion).Should(Equal("apps/v1"))
	Expect(sb.Spec.Workload.Kind).Should(Equal(workloadKind))
	Expect(sb.Spec.Workload.Name).Should(Equal(workloadName))
}
