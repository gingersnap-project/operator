package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KindEagerCacheRule = "EagerCacheRule"

// +kubebuilder:validation:Enum=Ready
type EagerCacheRuleConditionType string

const (
	EagerCacheRuleConditionReady EagerCacheRuleConditionType = "Ready"
)

// EagerCacheRuleCondition indicates the current status of a deployment
type EagerCacheRuleCondition struct {
	// Type is the type of the condition.
	Type EagerCacheRuleConditionType `json:"type,omitempty"`
	// +kubebuilder:validation:Enum=True;False;Unknown
	// Status is the status of the condition.
	Status metav1.ConditionStatus `json:"status,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// EagerCacheRuleStatus defines the observed state of EagerCacheRule
type EagerCacheRuleStatus struct {
	// +optional
	Conditions []EagerCacheRuleCondition `json:"conditions,omitempty"`
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
