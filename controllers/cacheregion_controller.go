package controllers

import (
	"context"
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"github.com/gingersnap-project/operator/pkg/reconcile/region"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gingersnapv1alpha1 "github.com/gingersnap-project/operator/api/v1alpha1"
)

// CacheRegionReconciler reconciles a CacheRegion object
type CacheRegionReconciler struct {
	*Reconciler
}

//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=cacheregions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=cacheregions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=cacheregions/finalizers,verbs=update

// Reconcile CacheRegion resources
func (r *CacheRegionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)

	instance := &v1alpha1.CacheRegion{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("CacheRegion CR not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, fmt.Errorf("unable to fetch CacheRegion CR %w", err)
	}

	// Don't reconcile CacheRegion CRs marked for deletion
	if instance.GetDeletionTimestamp() != nil {
		reqLogger.Info("Ignoring CacheRegion marked for deletion", "CacheRegion", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

	ctxProvider := reconcile.ContextProviderFunc(func(i interface{}) (reconcile.Context, error) {
		return r.NewPipelineCtx(ctx, reqLogger, instance), nil
	})
	retry, delay, err := region.PipelineBuilder().
		WithContextProvider(ctxProvider).
		Build().
		Process(instance)

	reqLogger.Info("Done", "requeue", retry, "requeueAfter", delay, "error", err)
	return ctrl.Result{Requeue: retry, RequeueAfter: delay}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *CacheRegionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gingersnapv1alpha1.CacheRegion{}).
		Complete(r)
}
