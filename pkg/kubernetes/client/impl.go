package client

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	metav1apply "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const fieldManager = "infinispan-operator"

var _ Client = &Runtime{}

// Runtime is a Client implementation based upon the controller-runtime client
type Runtime struct {
	record.EventRecorder
	Client    runtimeClient.Client
	Ctx       context.Context
	Namespace string
	Owner     runtimeClient.Object
	Scheme    *runtime.Scheme
}

func (c *Runtime) Apply(obj interface{}) error {
	// First convert to unstructured so that default values are emitted from the struct
	unstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}

	patch := &unstructured.Unstructured{
		Object: unstr,
	}

	patchOptions := &runtimeClient.PatchOptions{Force: pointer.Bool(true), FieldManager: fieldManager}
	if err := c.Client.Patch(c.Ctx, patch, runtimeClient.Apply, patchOptions); err != nil {
		return err
	}
	return nil
}

func (c *Runtime) OwnerReference() *metav1apply.OwnerReferenceApplyConfiguration {
	return OwnerReference(c.Owner)
}

func (c *Runtime) For(owner runtimeClient.Object) Client {
	clone := c.clone()
	clone.Owner = owner
	return clone
}

func (c *Runtime) WithNamespace(namespace string) Client {
	clone := c.clone()
	clone.Namespace = namespace
	return clone
}

func (c *Runtime) clone() *Runtime {
	return &Runtime{
		Ctx:           c.Ctx,
		Client:        c.Client,
		EventRecorder: c.EventRecorder,
		Namespace:     c.Namespace,
		Owner:         c.Owner,
		Scheme:        c.Scheme,
	}
}

func (c *Runtime) Create(obj runtimeClient.Object) error {
	return c.Client.Create(c.Ctx, obj)
}

func (c *Runtime) DeleteAllOf(set map[string]string, obj runtimeClient.Object, opts ...func(config *Config)) error {
	config := config(opts...)
	labelSelector := labels.SelectorFromSet(set)
	listOps := runtimeClient.ListOptions{LabelSelector: labelSelector}

	if !config.ClusterScoped() {
		listOps.Namespace = c.Namespace
	}
	deleteOpts := &runtimeClient.DeleteAllOfOptions{
		ListOptions: listOps,
	}
	return c.Client.DeleteAllOf(c.Ctx, obj, deleteOpts)
}

func (c *Runtime) Delete(name string, obj runtimeClient.Object, opts ...func(config *Config)) error {
	config := config(opts...)

	obj.SetName(name)
	if !config.ClusterScoped() {
		obj.SetNamespace(c.Namespace)
	}
	return c.Client.Delete(c.Ctx, obj)
}

func (c *Runtime) List(set map[string]string, list runtimeClient.ObjectList, opts ...func(config *Config)) error {
	config := config(opts...)
	labelSelector := labels.SelectorFromSet(set)
	listOps := &runtimeClient.ListOptions{LabelSelector: labelSelector}

	if !config.ClusterScoped() {
		listOps.Namespace = c.Namespace
	}
	return c.Client.List(c.Ctx, list, listOps)
}

func (c *Runtime) Load(name string, obj runtimeClient.Object, opts ...func(config *Config)) error {
	config := config(opts...)

	key := types.NamespacedName{
		Name: name,
	}
	if !config.ClusterScoped() {
		key.Namespace = c.Namespace
	}
	return c.Client.Get(c.Ctx, key, obj)
}

func (c *Runtime) Update(obj runtimeClient.Object) error {
	return c.Client.Update(c.Ctx, obj)
}

func (c *Runtime) UpdateStatus(obj runtimeClient.Object) error {
	return c.Client.Status().Update(c.Ctx, obj)
}

func OwnerReference(owner runtimeClient.Object) *metav1apply.OwnerReferenceApplyConfiguration {
	gvk := owner.GetObjectKind().GroupVersionKind()
	return metav1apply.OwnerReference().
		WithAPIVersion(gvk.GroupVersion().String()).
		WithKind(gvk.Kind).
		WithName(owner.GetName()).
		WithUID(owner.GetUID()).
		WithBlockOwnerDeletion(true).
		WithController(true)
}

func config(opts ...func(config *Config)) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}
