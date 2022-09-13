package cache

import (
	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/engytita/engytita-operator/pkg/reconcile/cache/context"
	"github.com/engytita/engytita-operator/pkg/reconcile/cache/infinispan"
	"github.com/engytita/engytita-operator/pkg/reconcile/cache/redis"
	"github.com/engytita/engytita-operator/pkg/reconcile/pipeline"
)

type HandlerFunc func(cache *v1alpha1.Cache, ctx *context.Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.Cache), ctx.(*context.Context))
}

func NewContextProvider(ctx reconcile.Context) reconcile.ContextProviderFunc {
	return func(i interface{}) (reconcile.Context, error) {
		return &context.Context{
			Context: ctx,
		}, nil
	}
}

func PipelineBuilder(cache *v1alpha1.Cache) *pipeline.Builder {
	builder := &pipeline.Builder{}
	if cache.Spec.Redis != nil {
		builder.WithHandlers(
			HandlerFunc(redis.Service),
			HandlerFunc(redis.ConfigurationSecret),
			HandlerFunc(redis.DaemonSet),
		)
	} else {
		builder.WithHandlers(
			HandlerFunc(infinispan.ConfigMap),
			HandlerFunc(infinispan.Service),
			HandlerFunc(infinispan.ConfigurationSecret),
			HandlerFunc(infinispan.DaemonSet),
		)
	}
	return builder
}
