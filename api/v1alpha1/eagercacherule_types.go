package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EagerCacheRuleSpec defines the desired state of EagerCacheRule
type EagerCacheRuleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of EagerCacheRule. Edit eagercacherule_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// EagerCacheRuleStatus defines the observed state of EagerCacheRule
type EagerCacheRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EagerCacheRule is the Schema for the eagercacherules API
type EagerCacheRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EagerCacheRuleSpec   `json:"spec,omitempty"`
	Status EagerCacheRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EagerCacheRuleList contains a list of EagerCacheRule
type EagerCacheRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EagerCacheRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EagerCacheRule{}, &EagerCacheRuleList{})
}
