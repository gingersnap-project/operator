package eager

import (
	"fmt"
	"time"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	binding "github.com/gingersnap-project/operator/pkg/apis/binding/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const conditionWait = time.Second * 2

func ConditionReady(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	ruleCondition := v1alpha1.EagerCacheRuleCondition{
		Type:   v1alpha1.EagerCacheRuleConditionReady,
		Status: metav1.ConditionFalse,
	}

	update := func(status metav1.ConditionStatus, msg string) {
		ruleCondition.Status = status
		ruleCondition.Message = msg
	}

	notFound := func(kind, name string) {
		ruleCondition.Status = metav1.ConditionFalse
		ruleCondition.Message = fmt.Sprintf("Cache %s '%s' %s", kind, name, metav1.StatusReasonNotFound)
	}

	cacheCondition := ctx.Cache.Condition(v1alpha1.CacheConditionReady)
	if cacheCondition.Status != metav1.ConditionTrue {
		ruleCondition.Message = fmt.Sprintf("Cache '%s' Not Ready", ctx.Cache.CacheService())
	} else {
		sb := &binding.ServiceBinding{}
		cache := ctx.Cache.CacheService()
		sbName := cache.DBSyncerCacheServiceBinding()
		if err := ctx.Client().Load(sbName, sb); err != nil {
			if errors.IsNotFound(err) {
				notFound("ServiceBinding", sbName)
			} else {
				ctx.RequeueAfter(conditionWait, fmt.Errorf("unable to load ServiceBinding '%s': %w", sbName, err))
				return
			}
		}

		var applicationBound bool
		for _, condition := range sb.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status == metav1.ConditionTrue {
					applicationBound = true
				} else {
					update(
						metav1.ConditionFalse,
						fmt.Sprintf("db-syncer ServiceBinding '%s' not Ready: '%s'", sbName, condition.Message),
					)
				}
				break
			}
		}

		if applicationBound {
			// Check that the db-syncer deployment is ready
			ready, msg, err := deploymentReady(ctx, cache.DBSyncerName())
			if err != nil {
				ctx.Requeue(err)
				return
			}

			if ready {
				update(
					metav1.ConditionTrue,
					"db-syncer Ready",
				)
			} else {
				update(
					metav1.ConditionFalse,
					msg,
				)
			}

			if r.Query() {
				// If Query is enabled for the rule, ensure that the index deployment is available
				ready, msg, err = deploymentReady(ctx, cache.SearchIndexName())
				if err != nil {
					ctx.Requeue(err)
					return
				}

				if ready {
					update(
						metav1.ConditionTrue,
						"All dependencies Ready",
					)
				} else {
					update(
						metav1.ConditionFalse,
						msg,
					)
				}
			}
		}
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

func deploymentReady(ctx *rule.Context, name string) (ready bool, msg string, err error) {
	deployment := &appsv1.Deployment{}
	if err = ctx.Client().Load(name, deployment); client.IgnoreNotFound(err) != nil {
		err = fmt.Errorf("unable to load Deployment '%s' for %s Ready Condition check: %w", name, v1alpha1.KindEagerCacheRule, err)
	} else if err != nil {
		msg = fmt.Sprintf("Cache Deployment '%s' %s", name, metav1.StatusReasonNotFound)
	} else {
		for _, condition := range deployment.Status.Conditions {
			if condition.Type == appsv1.DeploymentAvailable {
				if condition.Status == corev1.ConditionTrue {
					return true, "", nil
				}
			}
		}
		msg = fmt.Sprintf("Required '%s' Deployment '%d' pods to be Ready, observed '%d'", name, *deployment.Spec.Replicas, deployment.Status.ReadyReplicas)
	}
	return
}
