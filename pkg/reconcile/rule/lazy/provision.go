package lazy

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func LoadCache(r *v1alpha1.LazyCacheRule, ctx *rule.Context) {
	cacheRef := r.CacheService()
	cache := &v1alpha1.Cache{}
	err := ctx.Client().
		WithNamespace(cacheRef.Namespace).
		Load(cacheRef.Name, cache)

	if err != nil {
		msg := fmt.Sprintf("unable to load Cache CR '%s'", cacheRef)
		r.SetCondition(
			v1alpha1.LazyCacheRuleCondition{
				Type:    v1alpha1.LazyCacheRuleConditionReady,
				Status:  metav1.ConditionFalse,
				Message: msg,
			},
		)
		if err := ctx.Client().UpdateStatus(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to update Ready condition on LoadCache failure: %w", err))
			return
		}
		ctx.Requeue(fmt.Errorf("%s: %w", msg, err))
	}
	ctx.Cache = cache
}
