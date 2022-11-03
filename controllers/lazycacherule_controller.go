package controllers

import (
	"context"
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/pipeline"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule/lazy"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gingersnapv1alpha1 "github.com/gingersnap-project/operator/api/v1alpha1"
)

// LazyCacheRuleReconciler reconciles a LazyCacheRule object
type LazyCacheRuleReconciler struct {
	*Reconciler
}

//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=lazycacherules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=lazycacherules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gingersnap-project.io,namespace=gingersnap-operator-system,resources=lazycacherules/finalizers,verbs=update

// Reconcile LazyCacheRule resources
func (r *LazyCacheRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)

	instance := &v1alpha1.LazyCacheRule{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("LazyCacheRule CR not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, fmt.Errorf("unable to fetch LazyCacheRule CR %w", err)
	}

	var pipelineBuilder *pipeline.Builder
	if instance.GetDeletionTimestamp() != nil {
		pipelineBuilder = lazy.DeletePipelineBuilder()
	} else {
		pipelineBuilder = lazy.PipelineBuilder()
	}

	retry, delay, err := pipelineBuilder.
		WithContextProvider(
			rule.NewContextProvider(
				r.NewPipelineCtx(ctx, reqLogger, instance),
			),
		).
		Build().
		Process(instance)

	reqLogger.Info("Done", "requeue", retry, "requeueAfter", delay, "error", err)
	return ctrl.Result{Requeue: retry, RequeueAfter: delay}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *LazyCacheRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gingersnapv1alpha1.LazyCacheRule{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
