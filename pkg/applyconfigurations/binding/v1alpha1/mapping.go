// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// MappingApplyConfiguration represents an declarative configuration of the Mapping type for use
// with apply.
type MappingApplyConfiguration struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

// MappingApplyConfiguration constructs an declarative configuration of the Mapping type for use with
// apply.
func Mapping() *MappingApplyConfiguration {
	return &MappingApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *MappingApplyConfiguration) WithName(value string) *MappingApplyConfiguration {
	b.Name = &value
	return b
}

// WithValue sets the Value field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Value field is set to the value of the last call.
func (b *MappingApplyConfiguration) WithValue(value string) *MappingApplyConfiguration {
	b.Value = &value
	return b
}
