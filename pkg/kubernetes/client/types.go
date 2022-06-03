package client

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OperationResult is the result of a CreateOrPatch or CreateOrUpdate call
type OperationResult string

const ( // They should complete the sentence "Deployment default/foo has been ..."
	// OperationResultNone means that the resource has not been changed
	OperationResultNone OperationResult = "unchanged"
	// OperationResultCreated means that a new resource is created
	OperationResultCreated OperationResult = "created"
	// OperationResultUpdated means that an existing resource is updated
	OperationResultUpdated OperationResult = "updated"
	// OperationResultUpdatedStatus means that an existing resource and its status is updated
	OperationResultUpdatedStatus OperationResult = "updatedStatus"
	// OperationResultUpdatedStatusOnly means that only an existing status is updated
	OperationResultUpdatedStatusOnly OperationResult = "updatedStatusOnly"
)

type Config struct {
	setControllerRef bool
	clusterScoped    *bool
}

func (c *Config) ClusterScoped() bool {
	return c.clusterScoped != nil && *c.clusterScoped
}

func (c *Config) SetControllerRef() bool {
	return c.setControllerRef
}

// SetControllerRef indicates that the Owner object associated with the client should be set as this
// resources owner
func SetControllerRef(config *Config) {
	config.setControllerRef = true
}

// ClusterScoped indicates that the operation should be invoked on a cluster scoped resource
func ClusterScoped(config *Config) {
	config.clusterScoped = pointer.Bool(true)
}

type Client interface {
	record.EventRecorder
	// For returns a new Client implementation with the owner, used with SetControllerRef, set to the provided object
	For(owner v1.Object) Client
	// WithNamespace returns a new client implementation with the specified namespace
	WithNamespace(namespace string) Client
	// Create a k8s resource
	Create(obj client.Object, opts ...func(config *Config)) error
	// CreateOrUpdate an existing k8s resource using the mutate function., An OperationResult is returned
	// to indicate what the outcome of the operator was
	CreateOrUpdate(obj client.Object, mutate func() error, opts ...func(config *Config)) (OperationResult, error)
	// CreateOrPatch an existing k8s resource using the mutate function., An OperationResult is returned
	// to indicate what the outcome of the operator was.
	CreateOrPatch(obj client.Object, mutate func() error, opts ...func(config *Config)) (OperationResult, error)
	// Delete a k8s resource
	Delete(name string, obj client.Object, opts ...func(config *Config)) error
	// List k8s resources with labels matching those in the provided set
	List(set map[string]string, list client.ObjectList, opts ...func(config *Config)) error
	// Load a k8s resource
	Load(name string, obj client.Object, opts ...func(config *Config)) error
	// Update a k8s resource
	Update(obj client.Object) error
}
