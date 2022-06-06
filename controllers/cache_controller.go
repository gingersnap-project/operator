package controllers

import (
	"context"
	"fmt"

	"github.com/engytita/engytita-operator/api/v1alpha1"
	"github.com/engytita/engytita-operator/pkg/kubernetes/client"
	"github.com/engytita/engytita-operator/pkg/reconcile"
	"github.com/engytita/engytita-operator/pkg/reconcile/cache"
	"github.com/engytita/engytita-operator/pkg/reconcile/pipeline"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CacheReconciler reconciles a Cache object
type CacheReconciler struct {
	runtimeClient.Client
	Scheme *runtime.Scheme
	record.EventRecorder
}

//+kubebuilder:rbac:groups=engytita.org,resources=caches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=engytita.org,resources=caches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=engytita.org,resources=caches/finalizers,verbs=update

// Reconcile the Cache resource
func (r *CacheReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)

	instance := &v1alpha1.Cache{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Cache CR not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, fmt.Errorf("unable to fetch Cache CR %w", err)
	}

	// Don't reconcile Infinispan CRs marked for deletion
	if instance.GetDeletionTimestamp() != nil {
		reqLogger.Info(fmt.Sprintf("Ignoring Cache CR '%s:%s' marked for deletion", instance.Namespace, instance.Name))
		return ctrl.Result{}, nil
	}

	ctxProvider := reconcile.ContextProviderFunc(func(i interface{}) (reconcile.Context, error) {
		return pipeline.NewContext(ctx, reqLogger, &client.Runtime{
			Client:        r.Client,
			Ctx:           ctx,
			EventRecorder: r.EventRecorder,
			Scheme:        r.Scheme,
		}), nil
	})

	retry, delay, err := cache.Builder.
		WithContextProvider(ctxProvider).
		Build().
		Process(instance)

	reqLogger.Info("Done", "requeue", retry, "requeueAfter", delay, "error", err)
	return ctrl.Result{Requeue: retry, RequeueAfter: delay}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *CacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.Scheme = mgr.GetScheme()
	r.EventRecorder = mgr.GetEventRecorderFor("cache")
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Cache{}).
		Complete(r)
}
