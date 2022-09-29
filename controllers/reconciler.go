package controllers

import (
	"context"
	"fmt"

	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciler generic struct providing fields common to all reconciler structs
type Reconciler struct {
	runtimeClient.Client
	Scheme *runtime.Scheme
	record.EventRecorder
	supportedTypes map[schema.GroupVersionKind]struct{}
}

func (r *Reconciler) NewPipelineCtx(ctx context.Context, log logr.Logger, owner runtimeClient.Object) reconcile.Context {
	return pipeline.NewContext(ctx, log, r.supportedTypes, &client.Runtime{
		Client:        r.Client,
		Ctx:           ctx,
		EventRecorder: r.EventRecorder,
		Namespace:     owner.GetNamespace(),
		Owner:         owner,
		Scheme:        r.Scheme,
	})
}

func (r *Reconciler) InitSupportedTypes(mgr ctrl.Manager) error {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return fmt.Errorf("unable to create discovery client to determine supported types: %w", err)
	}

	types := []schema.GroupVersionKind{reconcile.ServiceMonitorGVK}
	supportedTypes := make(map[schema.GroupVersionKind]struct{}, len(types))
	for _, gvk := range types {
		groupVersion := gvk.GroupVersion().String()

		res, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion)
		if err == nil {
			for _, v := range res.APIResources {
				if v.Kind == gvk.Kind {
					supportedTypes[gvk] = struct{}{}
				}
			}
		} else if runtimeClient.IgnoreNotFound(err) != nil {
			return fmt.Errorf("unable to determine if '%s' is available: %w", gvk, err)
		}
	}
	r.supportedTypes = supportedTypes
	return nil
}
