// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// SafeTLSConfigApplyConfiguration represents an declarative configuration of the SafeTLSConfig type for use
// with apply.
type SafeTLSConfigApplyConfiguration struct {
	CA                 *SecretOrConfigMapApplyConfiguration `json:"ca,omitempty"`
	Cert               *SecretOrConfigMapApplyConfiguration `json:"cert,omitempty"`
	KeySecret          *corev1.SecretKeySelector            `json:"keySecret,omitempty"`
	ServerName         *string                              `json:"serverName,omitempty"`
	InsecureSkipVerify *bool                                `json:"insecureSkipVerify,omitempty"`
}

// SafeTLSConfigApplyConfiguration constructs an declarative configuration of the SafeTLSConfig type for use with
// apply.
func SafeTLSConfig() *SafeTLSConfigApplyConfiguration {
	return &SafeTLSConfigApplyConfiguration{}
}

// WithCA sets the CA field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CA field is set to the value of the last call.
func (b *SafeTLSConfigApplyConfiguration) WithCA(value *SecretOrConfigMapApplyConfiguration) *SafeTLSConfigApplyConfiguration {
	b.CA = value
	return b
}

// WithCert sets the Cert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Cert field is set to the value of the last call.
func (b *SafeTLSConfigApplyConfiguration) WithCert(value *SecretOrConfigMapApplyConfiguration) *SafeTLSConfigApplyConfiguration {
	b.Cert = value
	return b
}

// WithKeySecret sets the KeySecret field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the KeySecret field is set to the value of the last call.
func (b *SafeTLSConfigApplyConfiguration) WithKeySecret(value corev1.SecretKeySelector) *SafeTLSConfigApplyConfiguration {
	b.KeySecret = &value
	return b
}

// WithServerName sets the ServerName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServerName field is set to the value of the last call.
func (b *SafeTLSConfigApplyConfiguration) WithServerName(value string) *SafeTLSConfigApplyConfiguration {
	b.ServerName = &value
	return b
}

// WithInsecureSkipVerify sets the InsecureSkipVerify field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InsecureSkipVerify field is set to the value of the last call.
func (b *SafeTLSConfigApplyConfiguration) WithInsecureSkipVerify(value bool) *SafeTLSConfigApplyConfiguration {
	b.InsecureSkipVerify = &value
	return b
}
