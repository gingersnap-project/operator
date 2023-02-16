package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/cache"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CacheReconciler reconciles a Cache object
type CacheReconciler struct {
	*Reconciler
}

//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=caches,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=caches/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=caches/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,namespace=gingersnap-operator-system,resources=daemonsets,verbs=create;delete;deletecollection;get;list;patch;update;watch
// +kubebuilder:rbac:groups=apps,namespace=gingersnap-operator-system,resources=deployments,verbs=create;delete;get;list;patch;update;watch
// +kubebuilder:rbac:groups=core,namespace=gingersnap-operator-system,resources=secrets;services;configmaps,verbs=create;delete;deletecollection;get;list;patch;update;watch
// +kubebuilder:rbac:groups=core,namespace=gingersnap-operator-system,resources=serviceaccounts,verbs=create;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,namespace=gingersnap-operator-system,resources=roles;rolebindings,verbs=create;patch;

// +kubebuilder:rbac:groups=monitoring.coreos.com,namespace=gingersnap-operator-system,resources=servicemonitors,verbs=create;delete;get;list;patch;update;watch
// +kubebuilder:rbac:groups=servicebinding.io,namespace=gingersnap-operator-system,resources=servicebindings,verbs=create;get;list;patch;watch

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

	ctxProvider := cache.NewContextProvider(
		r.NewPipelineCtx(ctx, reqLogger, instance),
	)

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

	if err := r.InitSupportedTypes(mgr); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Cache{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
