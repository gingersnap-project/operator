package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (r *EagerCacheRule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gingersnap-project-io-v1alpha1-eagercacherule,mutating=true,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=eagercacherules,verbs=create;update,versions=v1alpha1,name=meagercacherule.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &EagerCacheRule{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *EagerCacheRule) Default() {
	if r.Spec.CacheRef == nil {
		// Do nothing as validation will fail
		return
	}

	if r.Spec.Key == nil {
		r.Spec.Key = &EagerCacheKey{}
	}
	r.Spec.Key.Format = KeyFormat_TEXT
	r.Spec.Key.KeySeparator = "|"
	r.CacheService().ApplyLabels(&r.ObjectMeta)
}

//+kubebuilder:webhook:path=/validate-gingersnap-project-io-v1alpha1-eagercacherule,mutating=false,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=eagercacherules,verbs=create;update,versions=v1alpha1,name=veagercacherule.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &EagerCacheRule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *EagerCacheRule) ValidateCreate() error {
	var allErrs field.ErrorList

	spec := field.NewPath("spec")

	if r.Spec.CacheRef == nil {
		allErrs = append(allErrs, field.Required(spec.Child("cacheRef"), "'cacheRef' field must be defined"))
	} else {
		cacheRefField := spec.Child("cacheRef")
		RequireField(&allErrs, "name", r.Spec.CacheRef.Name, cacheRefField)
		RequireField(&allErrs, "namespace", r.Spec.CacheRef.Name, cacheRefField)
	}

	RequireField(&allErrs, "tableName", r.Spec.TableName, spec)

	if r.Spec.Key == nil {
		allErrs = append(allErrs, field.Required(spec.Child("key"), FieldMustBeDefined("key")))
	} else {
		RequireNonEmptyArray(&allErrs, "keyColumns", r.Spec.Key.KeyColumns, spec.Child("key"))
	}
	return StatusError(allErrs, r.Name, KindEagerCacheRule)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *EagerCacheRule) ValidateUpdate(old runtime.Object) error {
	var allErrs field.ErrorList
	if err := EnsureRuleImmutability(&allErrs, KindEagerCacheRule, r, old.(*EagerCacheRule)); err != nil {
		return fmt.Errorf("unable to compare updated rule with existing rule: %w", err)
	}
	return StatusError(allErrs, r.Name, KindEagerCacheRule)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *EagerCacheRule) ValidateDelete() error {
	return nil
}
