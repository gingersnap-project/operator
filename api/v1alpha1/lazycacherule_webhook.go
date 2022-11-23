package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (r *LazyCacheRule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-gingersnap-project-io-v1alpha1-lazycacherule,mutating=true,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=lazycacherules,verbs=create;update,versions=v1alpha1,name=mlazycacherule.gingersnap-project.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &LazyCacheRule{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *LazyCacheRule) Default() {
	r.CacheService().ApplyLabels(&r.ObjectMeta)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-gingersnap-project-io-v1alpha1-lazycacherule,mutating=false,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=lazycacherules,verbs=create;update,versions=v1alpha1,name=vlazycacherule.gingersnap-project.io,admissionReviewVersions=v1

var _ webhook.Validator = &LazyCacheRule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LazyCacheRule) ValidateCreate() error {
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LazyCacheRule) ValidateUpdate(old runtime.Object) error {
	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LazyCacheRule) ValidateDelete() error {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
