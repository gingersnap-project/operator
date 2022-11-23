package eager

import (
	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
)

type HandlerFunc func(cache *v1alpha1.EagerCacheRule, ctx *rule.Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.EagerCacheRule), ctx.(*rule.Context))
}

func PipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	builder.WithHandlers(
		HandlerFunc(LoadCache),
		rule.HandlerFunc(rule.AddFinalizer),
		rule.HandlerFunc(rule.ApplyRuleConfigMap),
		HandlerFunc(ApplyDBSyncer),
	)
	return builder
}

func DeletePipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	return builder.
		WithHandlers(
			HandlerFunc(RemoveDBSyncer),
			rule.HandlerFunc(rule.RemoveRuleFromConfigMap),
			rule.HandlerFunc(rule.RemoveFinalizer),
		)
}
