package rule

import (
	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Context struct {
	reconcile.Context
	Cache *v1alpha1.Cache
}

type CacheRule interface {
	CacheService() v1alpha1.CacheService
	ConfigMap() string
	Filename() string
	Finalizer() string
	MarshallSpec() ([]byte, error)
	runtimeClient.Object
}

type HandlerFunc func(r CacheRule, ctx *Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(CacheRule), ctx.(*Context))
}

func NewContextProvider(ctx reconcile.Context) reconcile.ContextProviderFunc {
	return func(i interface{}) (reconcile.Context, error) {
		return &Context{
			Context: ctx,
		}, nil
	}
}
