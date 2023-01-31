package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KindCache = "Cache"

// +kubebuilder:validation:Enum=Ready
type CacheConditionType string

const (
	CacheConditionReady CacheConditionType = "Ready"
)

// CacheCondition indicates the current status of a deployment
type CacheCondition struct {
	// Type is the type of the condition.
	Type CacheConditionType `json:"type,omitempty"`
	// +kubebuilder:validation:Enum=True;False;Unknown
	// Status is the status of the condition.
	Status metav1.ConditionStatus `json:"status,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// CacheStatus defines the observed state of Cache
type CacheStatus struct {
	// +optional
	Conditions []CacheCondition `json:"conditions,omitempty"`
	// +optional
	ServiceBinding *ServiceBinding `json:"binding,omitempty"`
}

type ServiceBinding struct {
	Name string `json:"name,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cache is the Schema for the caches API
type Cache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CacheSpec   `json:"spec,omitempty"`
	Status CacheStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CacheList contains a list of Cache
type CacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cache{}, &CacheList{})
}
