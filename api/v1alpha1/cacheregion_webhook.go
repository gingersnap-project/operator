package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var cacheregionlog = logf.Log.WithName("cacheregion-resource")

func (r *CacheRegion) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-engytita-org-v1alpha1-cacheregion,mutating=true,failurePolicy=fail,sideEffects=None,groups=engytita.org,resources=cacheregions,verbs=create;update,versions=v1alpha1,name=mcacheregion.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &CacheRegion{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CacheRegion) Default() {
	cacheregionlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-engytita-org-v1alpha1-cacheregion,mutating=false,failurePolicy=fail,sideEffects=None,groups=engytita.org,resources=cacheregions,verbs=create;update,versions=v1alpha1,name=vcacheregion.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CacheRegion{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CacheRegion) ValidateCreate() error {
	cacheregionlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CacheRegion) ValidateUpdate(old runtime.Object) error {
	cacheregionlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CacheRegion) ValidateDelete() error {
	cacheregionlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
