package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/engytita/engytita-operator/api/v1alpha1"
	k8s "github.com/engytita/engytita-operator/pkg/kubernetes"
	"github.com/engytita/engytita-operator/pkg/reconcile/sidecar"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type PodReconciler struct {
	*Reconciler
}

// Reconcile Pods created with Engytita labels in order to set OwnerRef on created ConfigMap
func (p *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)

	pod := &corev1.Pod{}
	if err := p.Get(ctx, req.NamespacedName, pod); err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Pod not found", "pod", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, fmt.Errorf("unable to retrieve Pod %w", err)
	}

	// Don't reconcile Pods marked for deletion
	if pod.GetDeletionTimestamp() != nil {
		reqLogger.Info("Ignoring pod marked for deletion", "pod", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

	// Retrieve the name of the generated ConfigMap
	configMapVolume := k8s.Volume(sidecar.VolumeName, &pod.Spec)
	if configMapVolume == nil {
		return ctrl.Result{}, fmt.Errorf("unexpected state, sidecar volume '%s' not present", sidecar.VolumeName)
	}

	// Set this Pod as an Owner of the Region ConfigMap so that the ConfigMap is removed when no dependent Pods remain
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapVolume.ConfigMap.Name,
			Namespace: pod.Namespace,
		},
	}
	// Use CreateOrUpdate, instead of CreateOrPatch, to ensure any conflicts on ownerReference field are detected and
	_, err := controllerutil.CreateOrUpdate(ctx, p.Client, configMap, func() error {
		gvk := pod.GroupVersionKind()
		ref := metav1.OwnerReference{
			APIVersion:         gvk.GroupVersion().String(),
			Kind:               gvk.Kind,
			Name:               pod.Name,
			UID:                pod.UID,
			BlockOwnerDeletion: pointer.BoolPtr(true),
		}

		owners := configMap.GetOwnerReferences()
		if idx := indexOwnerRef(owners, ref); idx == -1 {
			owners = append(owners, ref)
		} else {
			owners[idx] = ref
		}
		configMap.SetOwnerReferences(owners)
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to set OwnerRef on region ConfigMap: %v", err)
	}
	return ctrl.Result{}, nil
}

// indexOwnerRef returns the index of the owner reference in the slice if found, or -1.
func indexOwnerRef(ownerReferences []metav1.OwnerReference, ref metav1.OwnerReference) int {
	for index, r := range ownerReferences {
		if referSameObject(r, ref) {
			return index
		}
	}
	return -1
}

// Returns true if a and b point to the same object.
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name && a.UID == b.UID
}

// SetupWithManager sets up the controller with the Manager.
func (p *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	p.Scheme = mgr.GetScheme()
	p.Client = mgr.GetClient()
	p.EventRecorder = mgr.GetEventRecorderFor(strings.ToLower(v1alpha1.KindCache))
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(event event.CreateEvent) bool {
				return v1alpha1.CacheServiceLabelsExist(event.Object.GetLabels())
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				return false
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				return false
			},
			GenericFunc: func(genericEvent event.GenericEvent) bool {
				return false
			},
		}).
		Complete(p)
}
