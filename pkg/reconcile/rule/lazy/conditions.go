package lazy

import (
	"fmt"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const conditionWait = time.Second * 2

func ConditionReady(r *v1alpha1.LazyCacheRule, ctx *rule.Context) {
	ruleCondition := v1alpha1.LazyCacheRuleCondition{
		Type:   v1alpha1.LazyCacheRuleConditionReady,
		Status: metav1.ConditionFalse,
	}
	cacheCondition := ctx.Cache.Condition(v1alpha1.CacheConditionReady)
	if cacheCondition.Status == metav1.ConditionTrue {
		ruleCondition.Status = metav1.ConditionTrue
		ruleCondition.Message = fmt.Sprintf("Cache '%s' Ready", ctx.Cache.CacheService())
	} else {
		ruleCondition.Message = fmt.Sprintf("Cache '%s' Not Ready", ctx.Cache.CacheService())
	}

	if r.SetCondition(ruleCondition) {
		if err := ctx.Client().UpdateStatus(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to update Ready condition: %w", err))
			return
		}
	}

	if ruleCondition.Status == metav1.ConditionFalse {
		ctx.RequeueAfter(conditionWait, nil)
	}
}
