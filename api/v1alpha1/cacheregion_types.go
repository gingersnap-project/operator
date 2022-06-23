package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CacheRegionSpec defines the desired state of CacheRegion
type CacheRegionSpec struct {
	// CacheReference defines the Cache that the CacheRegion is applied to
	Cache CacheService `json:"cache"`
}

// CacheService defines the location of the Cache resource that this CacheRegion should be applied to
type CacheService struct {
	// Name is the name of the Cache resource that the CacheRegion will be applied to
	Name string `json:"name"`
	// Namespace is the namespace in which the Cache CR belongs
	Namespace string `json:"namespace"`
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
