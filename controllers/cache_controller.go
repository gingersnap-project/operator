package controllers

import (
	"context"
	"fmt"
	"strings"

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

//+kubebuilder:rbac:groups=engytita.org,resources=caches,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=engytita.org,resources=caches/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=engytita.org,resources=caches/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,namespace=engytita-operator-system,resources=daemonsets,verbs=create;delete;deletecollection;get;list;patch;update;watch
// +kubebuilder:rbac:groups=core,namespace=engytita-operator-system,resources=services;configmaps,verbs=create;delete;deletecollection;get;list;patch;update;watch

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
			Namespace:     instance.Namespace,
			Owner:         instance,
			Scheme:        r.Scheme,
		}), nil
	})

	retry, delay, err := cache.PipelineBuilder(instance).
		WithContextProvider(ctxProvider).
		Build().
		Process(instance)

	reqLogger.Info("Done", "requeue", retry, "requeueAfter", delay, "error", err)
	return ctrl.Result{Requeue: retry, RequeueAfter: delay}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *CacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Scheme = mgr.GetScheme()
	r.Client = mgr.GetClient()
	r.EventRecorder = mgr.GetEventRecorderFor(strings.ToLower(v1alpha1.KindCache))
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Cache{}).
		Complete(r)
}
