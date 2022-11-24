package cache

import (
	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache/context"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache/infinispan"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
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

func PipelineBuilder(c *v1alpha1.Cache) *pipeline.Builder {
	builder := &pipeline.Builder{}

	var deploymentHandler HandlerFunc
	if c.Local() {
		deploymentHandler = infinispan.DaemonSet
	} else {
		deploymentHandler = infinispan.Deployment
	}

	return builder.WithHandlers(
		HandlerFunc(infinispan.WatchServiceAccount),
		HandlerFunc(infinispan.Service),
		HandlerFunc(infinispan.ConfigurationSecret),
		HandlerFunc(infinispan.ServiceMonitor),
		deploymentHandler,
	)
}
