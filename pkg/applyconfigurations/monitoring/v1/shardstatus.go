// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// ShardStatusApplyConfiguration represents an declarative configuration of the ShardStatus type for use
// with apply.
type ShardStatusApplyConfiguration struct {
	ShardID             *string `json:"shardID,omitempty"`
	Replicas            *int32  `json:"replicas,omitempty"`
	UpdatedReplicas     *int32  `json:"updatedReplicas,omitempty"`
	AvailableReplicas   *int32  `json:"availableReplicas,omitempty"`
	UnavailableReplicas *int32  `json:"unavailableReplicas,omitempty"`
}

// ShardStatusApplyConfiguration constructs an declarative configuration of the ShardStatus type for use with
// apply.
func ShardStatus() *ShardStatusApplyConfiguration {
	return &ShardStatusApplyConfiguration{}
}

// WithShardID sets the ShardID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ShardID field is set to the value of the last call.
func (b *ShardStatusApplyConfiguration) WithShardID(value string) *ShardStatusApplyConfiguration {
	b.ShardID = &value
	return b
}

// WithReplicas sets the Replicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Replicas field is set to the value of the last call.
func (b *ShardStatusApplyConfiguration) WithReplicas(value int32) *ShardStatusApplyConfiguration {
	b.Replicas = &value
	return b
}

// WithUpdatedReplicas sets the UpdatedReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UpdatedReplicas field is set to the value of the last call.
func (b *ShardStatusApplyConfiguration) WithUpdatedReplicas(value int32) *ShardStatusApplyConfiguration {
	b.UpdatedReplicas = &value
	return b
}

// WithAvailableReplicas sets the AvailableReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AvailableReplicas field is set to the value of the last call.
func (b *ShardStatusApplyConfiguration) WithAvailableReplicas(value int32) *ShardStatusApplyConfiguration {
	b.AvailableReplicas = &value
	return b
}

// WithUnavailableReplicas sets the UnavailableReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UnavailableReplicas field is set to the value of the last call.
func (b *ShardStatusApplyConfiguration) WithUnavailableReplicas(value int32) *ShardStatusApplyConfiguration {
	b.UnavailableReplicas = &value
	return b
}
