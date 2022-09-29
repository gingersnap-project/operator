// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

// MetadataConfigApplyConfiguration represents an declarative configuration of the MetadataConfig type for use
// with apply.
type MetadataConfigApplyConfiguration struct {
	Send         *bool        `json:"send,omitempty"`
	SendInterval *v1.Duration `json:"sendInterval,omitempty"`
}

// MetadataConfigApplyConfiguration constructs an declarative configuration of the MetadataConfig type for use with
// apply.
func MetadataConfig() *MetadataConfigApplyConfiguration {
	return &MetadataConfigApplyConfiguration{}
}

// WithSend sets the Send field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Send field is set to the value of the last call.
func (b *MetadataConfigApplyConfiguration) WithSend(value bool) *MetadataConfigApplyConfiguration {
	b.Send = &value
	return b
}

// WithSendInterval sets the SendInterval field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SendInterval field is set to the value of the last call.
func (b *MetadataConfigApplyConfiguration) WithSendInterval(value v1.Duration) *MetadataConfigApplyConfiguration {
	b.SendInterval = &value
	return b
}