// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// AlertmanagerConfigurationApplyConfiguration represents an declarative configuration of the AlertmanagerConfiguration type for use
// with apply.
type AlertmanagerConfigurationApplyConfiguration struct {
	Name *string `json:"name,omitempty"`
}

// AlertmanagerConfigurationApplyConfiguration constructs an declarative configuration of the AlertmanagerConfiguration type for use with
// apply.
func AlertmanagerConfiguration() *AlertmanagerConfigurationApplyConfiguration {
	return &AlertmanagerConfigurationApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *AlertmanagerConfigurationApplyConfiguration) WithName(value string) *AlertmanagerConfigurationApplyConfiguration {
	b.Name = &value
	return b
}
