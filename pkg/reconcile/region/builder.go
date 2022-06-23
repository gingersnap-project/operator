package region

import (
	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/engytita/engytita-operator/pkg/reconcile/pipeline"
)

type HandlerFunc func(cache *v1alpha1.CacheRegion, ctx reconcile.Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.CacheRegion), ctx)
}

func PipelineBuilder() *pipeline.Builder {
	builder := &pipeline.Builder{}
	return builder.WithHandlers(
		HandlerFunc(UpdateConfigMaps),
	)
}
