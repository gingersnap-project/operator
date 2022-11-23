package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (r *EagerCacheRule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-gingersnap-project-io-v1alpha1-eagercacherule,mutating=true,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=eagercacherules,verbs=create;update,versions=v1alpha1,name=meagercacherule.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &EagerCacheRule{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *EagerCacheRule) Default() {
	r.CacheService().ApplyLabels(&r.ObjectMeta)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-gingersnap-project-io-v1alpha1-eagercacherule,mutating=false,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=eagercacherules,verbs=create;update,versions=v1alpha1,name=veagercacherule.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &EagerCacheRule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *EagerCacheRule) ValidateCreate() error {
	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *EagerCacheRule) ValidateUpdate(old runtime.Object) error {
	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *EagerCacheRule) ValidateDelete() error {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
