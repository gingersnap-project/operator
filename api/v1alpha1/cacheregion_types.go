package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CacheRegionSpec defines the desired state of CacheRegion
type CacheRegionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of CacheRegion. Edit cacheregion_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// CacheRegionStatus defines the observed state of CacheRegion
type CacheRegionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CacheRegion is the Schema for the cacheregions API
type CacheRegion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CacheRegionSpec   `json:"spec,omitempty"`
	Status CacheRegionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CacheRegionList contains a list of CacheRegion
type CacheRegionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CacheRegion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CacheRegion{}, &CacheRegionList{})
}
