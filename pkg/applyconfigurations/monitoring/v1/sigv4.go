// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "k8s.io/api/core/v1"
)

// Sigv4ApplyConfiguration represents an declarative configuration of the Sigv4 type for use
// with apply.
type Sigv4ApplyConfiguration struct {
	Region    *string               `json:"region,omitempty"`
	AccessKey *v1.SecretKeySelector `json:"accessKey,omitempty"`
	SecretKey *v1.SecretKeySelector `json:"secretKey,omitempty"`
	Profile   *string               `json:"profile,omitempty"`
	RoleArn   *string               `json:"roleArn,omitempty"`
}

// Sigv4ApplyConfiguration constructs an declarative configuration of the Sigv4 type for use with
// apply.
func Sigv4() *Sigv4ApplyConfiguration {
	return &Sigv4ApplyConfiguration{}
}

// WithRegion sets the Region field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Region field is set to the value of the last call.
func (b *Sigv4ApplyConfiguration) WithRegion(value string) *Sigv4ApplyConfiguration {
	b.Region = &value
	return b
}

// WithAccessKey sets the AccessKey field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AccessKey field is set to the value of the last call.
func (b *Sigv4ApplyConfiguration) WithAccessKey(value v1.SecretKeySelector) *Sigv4ApplyConfiguration {
	b.AccessKey = &value
	return b
}

// WithSecretKey sets the SecretKey field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretKey field is set to the value of the last call.
func (b *Sigv4ApplyConfiguration) WithSecretKey(value v1.SecretKeySelector) *Sigv4ApplyConfiguration {
	b.SecretKey = &value
	return b
}

// WithProfile sets the Profile field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Profile field is set to the value of the last call.
func (b *Sigv4ApplyConfiguration) WithProfile(value string) *Sigv4ApplyConfiguration {
	b.Profile = &value
	return b
}

// WithRoleArn sets the RoleArn field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RoleArn field is set to the value of the last call.
func (b *Sigv4ApplyConfiguration) WithRoleArn(value string) *Sigv4ApplyConfiguration {
	b.RoleArn = &value
	return b
}
