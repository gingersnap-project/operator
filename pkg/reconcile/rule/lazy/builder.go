package lazy

import (
	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
)

type Context struct {
	reconcile.Context
	Cache *v1alpha1.Cache
}

func NewContextProvider(ctx reconcile.Context) reconcile.ContextProviderFunc {
	return func(i interface{}) (reconcile.Context, error) {
		return &Context{
			Context: ctx,
		}, nil
	}
}

type HandlerFunc func(cache *v1alpha1.LazyCacheRule, ctx *Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.LazyCacheRule), ctx.(*Context))
}

func PipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	builder.WithHandlers(
		HandlerFunc(LoadCache),
		HandlerFunc(AddFinalizer),
		HandlerFunc(AddRuleToConfigMap),
	)
	return builder
}

func DeletePipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	return builder.
		WithHandlers(
			HandlerFunc(RemoveRuleFromConfigMap),
			HandlerFunc(RemoveFinalizer),
		)
}
