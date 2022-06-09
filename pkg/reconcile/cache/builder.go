package cache

import (
	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/engytita/engytita-operator/pkg/reconcile/cache/infinispan"
	"github.com/engytita/engytita-operator/pkg/reconcile/cache/redis"
	"github.com/engytita/engytita-operator/pkg/reconcile/pipeline"
)

type HandlerFunc func(cache *v1alpha1.Cache, ctx reconcile.Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.Cache), ctx)
}

func PipelineBuilder(cache *v1alpha1.Cache) *pipeline.Builder {
	builder := &pipeline.Builder{}
	if cache.Spec.Redis != nil {
		builder.WithHandlers(
			HandlerFunc(redis.Service),
			HandlerFunc(redis.DaemonSet),
		)
	} else {
		builder.WithHandlers(
			HandlerFunc(infinispan.ConfigMap),
			HandlerFunc(infinispan.Service),
			HandlerFunc(infinispan.DaemonSet),
		)
	}
	return builder
}
