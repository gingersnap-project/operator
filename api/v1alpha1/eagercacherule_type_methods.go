package v1alpha1

import (
	"google.golang.org/protobuf/encoding/protojson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (r *EagerCacheRule) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
}

func (r *EagerCacheRule) Finalizer() string {
	return schema.GroupKind{Group: Group, Kind: KindEagerCacheRule}.String()
}

func (r *EagerCacheRule) CacheService() CacheService {
	return CacheService{
		Name:      r.Spec.CacheRef.Name,
		Namespace: r.Spec.CacheRef.Namespace,
	}
}

func (r *EagerCacheRule) ConfigMap() string {
	return r.CacheService().EagerCacheConfigMap()
}

func (r *EagerCacheRule) MarshallSpec() ([]byte, error) {
	return protojson.MarshalOptions{Multiline: true}.Marshal(&r.Spec)
}

func (r *EagerCacheRule) Condition(condition EagerCacheRuleConditionType) EagerCacheRuleCondition {
	for _, existing := range r.Status.Conditions {
		if existing.Type == condition {
			return existing
		}
	}
	// Absence of condition means `False` value
	return EagerCacheRuleCondition{Type: condition, Status: metav1.ConditionFalse}
}

func (r *EagerCacheRule) SetCondition(condition EagerCacheRuleCondition) (updated bool) {
	for idx := range r.Status.Conditions {
		c := &r.Status.Conditions[idx]
		if c.Type == condition.Type {
			if c.Status != condition.Status {
				c.Status = condition.Status
				updated = true
			}
			if c.Message != condition.Message {
				c.Message = condition.Message
				updated = true
			}
			return updated
		}
	}
	r.Status.Conditions = append(r.Status.Conditions, EagerCacheRuleCondition{
		Type:    condition.Type,
		Status:  condition.Status,
		Message: condition.Message,
	})
	return true
}
