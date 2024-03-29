package lazy

import (
	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
)

type HandlerFunc func(r *v1alpha1.LazyCacheRule, ctx *rule.Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.LazyCacheRule), ctx.(*rule.Context))
}

func PipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	builder.WithHandlers(
		HandlerFunc(LoadCache),
		rule.HandlerFunc(rule.AddFinalizer),
		rule.HandlerFunc(rule.ApplyRuleConfigMap),
		HandlerFunc(ConditionReady),
	)
	return builder
}

func DeletePipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	return builder.
		WithHandlers(
			rule.HandlerFunc(rule.RemoveRuleFromConfigMap),
			rule.HandlerFunc(rule.RemoveFinalizer),
		)
}
