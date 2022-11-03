package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KindLazyCacheRule = "LazyCacheRule"

// LazyCacheRuleSpec defines the desired state of LazyCacheRule
type LazyCacheRuleSpec struct {
	// CacheReference defines the Cache that the LazyCacheRule is applied to
	Cache CacheService `json:"cache"`
}

// LazyCacheRuleStatus defines the observed state of LazyCacheRule
type LazyCacheRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LazyCacheRule is the Schema for the lazycacherules API
type LazyCacheRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LazyCacheRuleSpec   `json:"spec,omitempty"`
	Status LazyCacheRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LazyCacheRuleList contains a list of LazyCacheRule
type LazyCacheRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LazyCacheRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LazyCacheRule{}, &LazyCacheRuleList{})
}
