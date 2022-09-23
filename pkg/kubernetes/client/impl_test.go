package client_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	k8sClient  runtimeClient.Client
	testEnv    *envtest.Environment
	ctx        context.Context
	cancel     context.CancelFunc
	namespace  = "default"
	testClient client.Client
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Kubernetes Client",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	scheme := runtime.NewScheme()
	Expect(corev1.SchemeBuilder.AddToScheme(scheme)).Should(Succeed())

	k8sClient, err = runtimeClient.New(cfg, runtimeClient.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	testClient = &client.Runtime{
		Client:        k8sClient,
		Ctx:           ctx,
		EventRecorder: nil,
		Namespace:     namespace,
		Scheme:        scheme,
	}
}, 60)

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Runtime", func() {

	AfterEach(func() {
		const timeout = time.Second * 30
		const interval = time.Second * 1

		// Delete created resources
		By("Expecting to delete successfully")
		deleteOpts := []runtimeClient.DeleteAllOfOption{
			runtimeClient.InNamespace(namespace),
		}
		Expect(k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, deleteOpts...)).Should(Succeed())
		Expect(k8sClient.DeleteAllOf(ctx, &corev1.Secret{}, deleteOpts...)).Should(Succeed())

		By("Expecting to delete finish")
		Eventually(func() int {
			cmList := &corev1.ConfigMapList{}
			secretList := &corev1.SecretList{}
			listOps := &runtimeClient.ListOptions{Namespace: namespace}
			Expect(k8sClient.List(ctx, cmList, listOps)).Should(Succeed())
			Expect(k8sClient.List(ctx, secretList, listOps)).Should(Succeed())
			return len(cmList.Items) + len(secretList.Items)
		}, timeout, interval).Should(Equal(0))
	})

	It("should create and load resources", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "test-cm",
			},
			Data: map[string]string{"key": "value"},
		}
		Expect(testClient.Create(cm)).Should(Succeed())

		created := &corev1.ConfigMap{}
		Expect(testClient.Load(cm.Name, created))
		Expect(created.Data["key"]).Should(Equal("value"))
		Expect(created.OwnerReferences).Should(BeEmpty())
	})

	It("should apply and load resources", func() {
		cm := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "test-cm",
			},
			Data: map[string]string{"key": "value"},
		}
		Expect(testClient.Apply(cm)).Should(Succeed())

		created := &corev1.ConfigMap{}
		Expect(testClient.Load(cm.Name, created))
		Expect(created.Data["key"]).Should(Equal("value"))
		Expect(created.OwnerReferences).Should(BeEmpty())
	})

	It("should load cluster scoped resources", func() {
		ns := &corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Namespace",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "some-namespace",
			},
		}
		Expect(testClient.Apply(ns)).Should(Succeed())

		created := &corev1.Namespace{}
		Expect(testClient.Load(ns.Name, created, client.ClusterScoped)).Should(Succeed())
		Expect(created.Name).Should(Equal(ns.Name))
	})
})
