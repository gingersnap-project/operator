package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (r *LazyCacheRule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gingersnap-project-io-v1alpha1-lazycacherule,mutating=true,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=lazycacherules,verbs=create;update,versions=v1alpha1,name=mlazycacherule.gingersnap-project.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &LazyCacheRule{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *LazyCacheRule) Default() {
	if r.Spec.CacheRef == nil {
		// Do nothing as validation will fail
		return
	}

	if r.Spec.Key == nil {
		r.Spec.Key = &LazyCacheKey{}
	}
	r.Spec.Key.Format = KeyFormat_TEXT
	r.Spec.Key.KeySeparator = "|"
	r.CacheService().ApplyLabels(&r.ObjectMeta)
}

//+kubebuilder:webhook:path=/validate-gingersnap-project-io-v1alpha1-lazycacherule,mutating=false,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=lazycacherules,verbs=create;update,versions=v1alpha1,name=vlazycacherule.gingersnap-project.io,admissionReviewVersions=v1

var _ webhook.Validator = &LazyCacheRule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LazyCacheRule) ValidateCreate() error {
	var allErrs field.ErrorList
	if r.Spec.CacheRef == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("cacheRef"), "'cacheRef' field must be defined"))
	} else {
		cacheRefField := field.NewPath("spec").Child("cacheRef")
		RequireField(&allErrs, "name", r.Spec.CacheRef.Name, cacheRefField)
		RequireField(&allErrs, "namespace", r.Spec.CacheRef.Name, cacheRefField)
	}
	RequireField(&allErrs, "query", r.Spec.Query, field.NewPath("spec"))
	return StatusError(allErrs, r.Name, KindLazyCacheRule)
}

func (r *LazyCacheRule) ValidateUpdate(old runtime.Object) error {
	var allErrs field.ErrorList
	if err := EnsureRuleImmutability(&allErrs, KindLazyCacheRule, r, old.(*LazyCacheRule)); err != nil {
		return fmt.Errorf("unable to compare updated rule with existing rule: %w", err)
	}
	return StatusError(allErrs, r.Name, KindLazyCacheRule)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LazyCacheRule) ValidateDelete() error {
	return nil
}
