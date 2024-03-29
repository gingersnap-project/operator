// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// EmbeddedObjectMetadataApplyConfiguration represents an declarative configuration of the EmbeddedObjectMetadata type for use
// with apply.
type EmbeddedObjectMetadataApplyConfiguration struct {
	Name        *string           `json:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// EmbeddedObjectMetadataApplyConfiguration constructs an declarative configuration of the EmbeddedObjectMetadata type for use with
// apply.
func EmbeddedObjectMetadata() *EmbeddedObjectMetadataApplyConfiguration {
	return &EmbeddedObjectMetadataApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *EmbeddedObjectMetadataApplyConfiguration) WithName(value string) *EmbeddedObjectMetadataApplyConfiguration {
	b.Name = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *EmbeddedObjectMetadataApplyConfiguration) WithLabels(entries map[string]string) *EmbeddedObjectMetadataApplyConfiguration {
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *EmbeddedObjectMetadataApplyConfiguration) WithAnnotations(entries map[string]string) *EmbeddedObjectMetadataApplyConfiguration {
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}
