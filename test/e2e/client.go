package e2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/engytita/engytita-operator/pkg/kubernetes/client"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TestClient struct {
	client.Client
	Ctx  context.Context
	Rest rest.Interface
}

func (c *TestClient) WriteAllResourcesToFile(dir string) {
	// Operator Pod logs
	printErr(os.MkdirAll(dir, os.ModePerm))
	c.WriteKindToFile(dir, OperatorNamespace, "Pod", &corev1.PodList{}, map[string]string{"app.kubernetes.io/name": "engytita-operator"})
}

func (c *TestClient) WriteKindToFile(dir, namespace, suffix string, list runtimeClient.ObjectList, set labels.Set) {
	err := c.WithNamespace(namespace).List(set, list)
	printErr(err)

	if podList, ok := list.(*corev1.PodList); ok {
		for _, pod := range podList.Items {
			for _, container := range pod.Spec.Containers {
				log, err := c.Logs(pod.Name, container.Name, namespace)
				printErr(err)

				fileName := fmt.Sprintf("%s/%s_%s.log", dir, pod.Name, container.Name)
				err = ioutil.WriteFile(fileName, []byte(log), 0666)
				printErr(err)
			}
		}
	}

	unstructuredResource, err := runtime.DefaultUnstructuredConverter.ToUnstructured(list)
	printErr(err)
	unstructuredResourceList := unstructured.UnstructuredList{}
	unstructuredResourceList.SetUnstructuredContent(unstructuredResource)

	for _, item := range unstructuredResourceList.Items {
		yaml_, err := yaml.Marshal(item)
		printErr(err)

		fileName := fmt.Sprintf("%s/%s_%s.yaml", dir, item.GetName(), suffix)
		err = ioutil.WriteFile(fileName, yaml_, 0666)
		printErr(err)
	}
}

func (c *TestClient) Logs(pod, container, namespace string) (string, error) {
	req := c.Rest.Get().Namespace(namespace).Resource("pods").Name(pod).SubResource("log")
	if container != "" {
		req.Param("container", container)
	}

	readCloser, err := req.Stream(c.Ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		cerr := readCloser.Close()
		if err == nil {
			err = cerr
		}
	}()

	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	return string(body), err
}

// printErr if the error is not nil, printErr to stdout. Should only be used by test cleanup operations that shouldn't result
// in a test failing
func printErr(err error) {
	if err != nil {
		fmt.Printf("Encountered error: %v", err)
	}
}
