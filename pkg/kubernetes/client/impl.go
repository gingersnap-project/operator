package client

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Runtime is a Client implementation based upon the controller-runtime client
type Runtime struct {
	record.EventRecorder
	Client    runtimeClient.Client
	Ctx       context.Context
	Namespace string
	Owner     metav1.Object
	Scheme    *runtime.Scheme
}

func (c *Runtime) For(owner metav1.Object) Client {
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

func (c *Runtime) Create(obj runtimeClient.Object, opts ...func(config *Config)) error {
	if err := c.setControllerRef(obj, config(opts...)); err != nil {
		return err
	}
	return c.Client.Create(c.Ctx, obj)
}

func (c *Runtime) CreateOrUpdate(obj runtimeClient.Object, mutate func() error, opts ...func(config *Config)) (OperationResult, error) {
	res, err := controllerutil.CreateOrUpdate(c.Ctx, c.Client, obj, func() error {
		if mutate != nil {
			if err := mutate(); err != nil {
				return err
			}
		}
		return c.setControllerRef(obj, config(opts...))
	})
	return OperationResult(res), err
}

func (c *Runtime) CreateOrPatch(obj runtimeClient.Object, mutate func() error, opts ...func(config *Config)) (OperationResult, error) {
	res, err := controllerutil.CreateOrPatch(c.Ctx, c.Client, obj, func() error {
		if mutate != nil {
			if err := mutate(); err != nil {
				return err
			}
		}
		return c.setControllerRef(obj, config(opts...))
	})
	return OperationResult(res), err
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

func config(opts ...func(config *Config)) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func (c *Runtime) setControllerRef(obj runtimeClient.Object, config *Config) error {
	if config.setControllerRef {
		if c.Owner == nil {
			return fmt.Errorf("unable to SetControllerRef, Owner is nil")
		}
		if err := controllerutil.SetControllerReference(c.Owner, obj, c.Scheme); err != nil {
			return err
		}
	}
	return nil
}
