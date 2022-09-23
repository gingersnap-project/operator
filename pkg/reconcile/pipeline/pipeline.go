package pipeline

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gingersnap-project/operator/pkg/reconcile"
)

var _ reconcile.Pipeline = &impl{}

type impl struct {
	ctxProvider reconcile.ContextProvider
	handlers    []reconcile.Handler
}

func (i *impl) Process(resource interface{}) (retry bool, delay time.Duration, err error) {
	defer func() {
		if perr := recover(); perr != nil {
			retry = true
			err = fmt.Errorf("panic occurred: %v", perr)
		}
	}()
	context, err := i.ctxProvider.Get(resource)
	if err != nil {
		return false, 0, err
	}

	var status reconcile.FlowStatus
	for _, h := range i.handlers {
		invokeHandler(h, resource, context)
		status = context.Status()
		if status.Stop {
			break
		}
	}
	return status.Retry, status.Delay, status.Err
}

func invokeHandler(h reconcile.Handler, i interface{}, ctx reconcile.Context) {
	defer func() {
		if err := recover(); err != nil {
			e := fmt.Errorf("panic occurred: %v", err)
			ctx.Log().Error(e, string(debug.Stack()))
			ctx.Requeue(e)
		}
	}()
	//fmt.Printf("Handler=%s\n", runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
	h.Handle(i, ctx)
}
