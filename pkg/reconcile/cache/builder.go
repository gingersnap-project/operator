package cache

import (
	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
)

type Context struct {
	reconcile.Context
}

type HandlerFunc func(cache *v1alpha1.Cache, ctx *Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.Cache), ctx.(*Context))
}

func NewContextProvider(ctx reconcile.Context) reconcile.ContextProviderFunc {
	return func(i interface{}) (reconcile.Context, error) {
		return &Context{
			Context: ctx,
		}, nil
	}
}

func PipelineBuilder(c *v1alpha1.Cache) *pipeline.Builder {
	builder := &pipeline.Builder{}

	var deploymentHandler HandlerFunc
	if c.Local() {
		deploymentHandler = DaemonSet
	} else {
		deploymentHandler = Deployment
	}

	return builder.WithHandlers(
		HandlerFunc(WatchServiceAccount),
		HandlerFunc(Service),
		HandlerFunc(UserServiceBindingSecret),
		HandlerFunc(DBSyncerCacheServiceBindingSecret),
		HandlerFunc(ApplyDataSourceServiceBinding),
		HandlerFunc(ServiceMonitor),
		deploymentHandler,
		HandlerFunc(ConditionReady),
	)
}
