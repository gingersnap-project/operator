// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// QueueConfigApplyConfiguration represents an declarative configuration of the QueueConfig type for use
// with apply.
type QueueConfigApplyConfiguration struct {
	Capacity          *int    `json:"capacity,omitempty"`
	MinShards         *int    `json:"minShards,omitempty"`
	MaxShards         *int    `json:"maxShards,omitempty"`
	MaxSamplesPerSend *int    `json:"maxSamplesPerSend,omitempty"`
	BatchSendDeadline *string `json:"batchSendDeadline,omitempty"`
	MaxRetries        *int    `json:"maxRetries,omitempty"`
	MinBackoff        *string `json:"minBackoff,omitempty"`
	MaxBackoff        *string `json:"maxBackoff,omitempty"`
	RetryOnRateLimit  *bool   `json:"retryOnRateLimit,omitempty"`
}

// QueueConfigApplyConfiguration constructs an declarative configuration of the QueueConfig type for use with
// apply.
func QueueConfig() *QueueConfigApplyConfiguration {
	return &QueueConfigApplyConfiguration{}
}

// WithCapacity sets the Capacity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Capacity field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithCapacity(value int) *QueueConfigApplyConfiguration {
	b.Capacity = &value
	return b
}

// WithMinShards sets the MinShards field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinShards field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithMinShards(value int) *QueueConfigApplyConfiguration {
	b.MinShards = &value
	return b
}

// WithMaxShards sets the MaxShards field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxShards field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithMaxShards(value int) *QueueConfigApplyConfiguration {
	b.MaxShards = &value
	return b
}

// WithMaxSamplesPerSend sets the MaxSamplesPerSend field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxSamplesPerSend field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithMaxSamplesPerSend(value int) *QueueConfigApplyConfiguration {
	b.MaxSamplesPerSend = &value
	return b
}

// WithBatchSendDeadline sets the BatchSendDeadline field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BatchSendDeadline field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithBatchSendDeadline(value string) *QueueConfigApplyConfiguration {
	b.BatchSendDeadline = &value
	return b
}

// WithMaxRetries sets the MaxRetries field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxRetries field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithMaxRetries(value int) *QueueConfigApplyConfiguration {
	b.MaxRetries = &value
	return b
}

// WithMinBackoff sets the MinBackoff field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinBackoff field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithMinBackoff(value string) *QueueConfigApplyConfiguration {
	b.MinBackoff = &value
	return b
}

// WithMaxBackoff sets the MaxBackoff field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxBackoff field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithMaxBackoff(value string) *QueueConfigApplyConfiguration {
	b.MaxBackoff = &value
	return b
}

// WithRetryOnRateLimit sets the RetryOnRateLimit field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RetryOnRateLimit field is set to the value of the last call.
func (b *QueueConfigApplyConfiguration) WithRetryOnRateLimit(value bool) *QueueConfigApplyConfiguration {
	b.RetryOnRateLimit = &value
	return b
}
