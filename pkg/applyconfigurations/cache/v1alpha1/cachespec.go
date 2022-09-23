// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/gingersnap-project/operator/api/v1alpha1"
)

// CacheSpecApplyConfiguration represents an declarative configuration of the CacheSpec type for use
// with apply.
type CacheSpecApplyConfiguration struct {
	Infinispan *v1alpha1.InfinispanSpec `json:"infinispan,omitempty"`
	Redis      *v1alpha1.RedisSpec      `json:"redis,omitempty"`
}

// CacheSpecApplyConfiguration constructs an declarative configuration of the CacheSpec type for use with
// apply.
func CacheSpec() *CacheSpecApplyConfiguration {
	return &CacheSpecApplyConfiguration{}
}

// WithInfinispan sets the Infinispan field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Infinispan field is set to the value of the last call.
func (b *CacheSpecApplyConfiguration) WithInfinispan(value v1alpha1.InfinispanSpec) *CacheSpecApplyConfiguration {
	b.Infinispan = &value
	return b
}

// WithRedis sets the Redis field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Redis field is set to the value of the last call.
func (b *CacheSpecApplyConfiguration) WithRedis(value v1alpha1.RedisSpec) *CacheSpecApplyConfiguration {
	b.Redis = &value
	return b
}
