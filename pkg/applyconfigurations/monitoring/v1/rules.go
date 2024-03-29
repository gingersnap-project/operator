// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// RulesApplyConfiguration represents an declarative configuration of the Rules type for use
// with apply.
type RulesApplyConfiguration struct {
	Alert *RulesAlertApplyConfiguration `json:"alert,omitempty"`
}

// RulesApplyConfiguration constructs an declarative configuration of the Rules type for use with
// apply.
func Rules() *RulesApplyConfiguration {
	return &RulesApplyConfiguration{}
}

// WithAlert sets the Alert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Alert field is set to the value of the last call.
func (b *RulesApplyConfiguration) WithAlert(value *RulesAlertApplyConfiguration) *RulesApplyConfiguration {
	b.Alert = value
	return b
}
