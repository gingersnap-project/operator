package controllers

import (
	"context"

	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciler generic struct providing fields common to all reconciler structs
type Reconciler struct {
	runtimeClient.Client
	Scheme *runtime.Scheme
	record.EventRecorder
}

func (r *Reconciler) NewPipelineCtx(ctx context.Context, log logr.Logger, owner runtimeClient.Object) reconcile.Context {
	return pipeline.NewContext(ctx, log, &client.Runtime{
		Client:        r.Client,
		Ctx:           ctx,
		EventRecorder: r.EventRecorder,
		Namespace:     owner.GetNamespace(),
		Owner:         owner,
		Scheme:        r.Scheme,
	})
}
