// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// RefApplyConfiguration represents an declarative configuration of the Ref type for use
// with apply.
type RefApplyConfiguration struct {
	Group    *string `json:"group,omitempty"`
	Version  *string `json:"version,omitempty"`
	Kind     *string `json:"kind,omitempty"`
	Resource *string `json:"resource,omitempty"`
	Name     *string `json:"name,omitempty"`
}

// RefApplyConfiguration constructs an declarative configuration of the Ref type for use with
// apply.
func Ref() *RefApplyConfiguration {
	return &RefApplyConfiguration{}
}

// WithGroup sets the Group field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Group field is set to the value of the last call.
func (b *RefApplyConfiguration) WithGroup(value string) *RefApplyConfiguration {
	b.Group = &value
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *RefApplyConfiguration) WithVersion(value string) *RefApplyConfiguration {
	b.Version = &value
	return b
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *RefApplyConfiguration) WithKind(value string) *RefApplyConfiguration {
	b.Kind = &value
	return b
}

// WithResource sets the Resource field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Resource field is set to the value of the last call.
func (b *RefApplyConfiguration) WithResource(value string) *RefApplyConfiguration {
	b.Resource = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *RefApplyConfiguration) WithName(value string) *RefApplyConfiguration {
	b.Name = &value
	return b
}
