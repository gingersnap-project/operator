//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	k8sClient *TestClient
	testEnv   *envtest.Environment
	cancel    context.CancelFunc
	ctx       context.Context
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)

	if err := os.RemoveAll(OutputDir); err != nil {
		t.Errorf("Failed removing report directory: %v", err)
	}

	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		t.Errorf("Failed creating report directory: %v", err)
	}

	fmt.Printf("Writing test output files to: '%s'\n", OutputDir)
	fmt.Printf("Operator Namespace: '%s'\n", OperatorNamespace)
	fmt.Printf("Test Namespace: '%s'\n", Namespace)
	RunSpecsWithDefaultAndCustomReporters(t,
		"E2E",
		[]Reporter{
			printer.NewlineReporter{},
			reporters.NewJUnitReporter(filepath.Join(OutputDir, "junit.xml")),
		},
	)
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		UseExistingCluster: pointer.Bool(true),
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	scheme := runtime.NewScheme()
	Expect(corev1.SchemeBuilder.AddToScheme(scheme)).Should(Succeed())
	Expect(appsv1.SchemeBuilder.AddToScheme(scheme)).Should(Succeed())
	Expect(admissionv1beta1.AddToScheme(scheme)).Should(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).Should(Succeed())

	runtime, err := runtimeClient.New(cfg, runtimeClient.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(runtime).NotTo(BeNil())

	config := testEnv.Config
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.APIPath = "/api"
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	restClient, err := rest.RESTClientFor(config)
	Expect(err).NotTo(HaveOccurred())

	k8sClient = &TestClient{
		Client: &client.Runtime{
			Client:    runtime,
			Ctx:       ctx,
			Namespace: Namespace,
			Scheme:    scheme,
		},
		Ctx:  ctx,
		Rest: restClient,
	}

	if MultiNamespace {
		Expect(k8sClient.Create(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: Namespace}})).Should(Succeed())
	}
}, 60)

var _ = AfterSuite(func() {
	if CleanupTestNamespace && MultiNamespace {
		Expect(k8sClient.Delete(Namespace, &corev1.Namespace{}, client.ClusterScoped)).Should(Succeed())
	}
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
