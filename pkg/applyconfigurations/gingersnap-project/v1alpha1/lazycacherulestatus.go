// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// LazyCacheRuleStatusApplyConfiguration represents an declarative configuration of the LazyCacheRuleStatus type for use
// with apply.
type LazyCacheRuleStatusApplyConfiguration struct {
	Conditions []LazyCacheRuleConditionApplyConfiguration `json:"conditions,omitempty"`
}

// LazyCacheRuleStatusApplyConfiguration constructs an declarative configuration of the LazyCacheRuleStatus type for use with
// apply.
func LazyCacheRuleStatus() *LazyCacheRuleStatusApplyConfiguration {
	return &LazyCacheRuleStatusApplyConfiguration{}
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
func (b *LazyCacheRuleStatusApplyConfiguration) WithConditions(values ...*LazyCacheRuleConditionApplyConfiguration) *LazyCacheRuleStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}
