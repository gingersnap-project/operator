package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KindLazyCacheRule = "LazyCacheRule"

// +kubebuilder:validation:Enum=Ready
type LazyCacheRuleConditionType string

const (
	LazyCacheRuleConditionReady LazyCacheRuleConditionType = "Ready"
)

// LazyCacheRuleCondition indicates the current status of a deployment
type LazyCacheRuleCondition struct {
	// Type is the type of the condition.
	Type LazyCacheRuleConditionType `json:"type,omitempty"`
	// +kubebuilder:validation:Enum=True;False;Unknown
	// Status is the status of the condition.
	Status metav1.ConditionStatus `json:"status,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// LazyCacheRuleStatus defines the observed state of LazyCacheRule
type LazyCacheRuleStatus struct {
	// +optional
	Conditions []LazyCacheRuleCondition `json:"conditions,omitempty"`
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
