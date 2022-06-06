package cache

import (
	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/engytita/engytita-operator/pkg/reconcile/pipeline"
)

var Builder = pipeline.Builder().WithHandlers(defaultFlow...)

var defaultFlow = []reconcile.Handler{
	HandlerFunc(Example),
}

type HandlerFunc func(cache *v1alpha1.Cache, ctx reconcile.Context)

func (f HandlerFunc) Handle(i interface{}, ctx reconcile.Context) {
	f(i.(*v1alpha1.Cache), ctx)
}

// TODO add handler implementations
func Example(cache *v1alpha1.Cache, ctx reconcile.Context) {

}
