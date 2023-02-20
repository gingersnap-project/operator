package v1alpha1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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

// RegisterEagerRuleValidatingWebhook explicitly adds the validating webhook to the Webhook Server
// This is necessary as we need to implement admission.Handler interface directly so that the request context can be
// used by the runtime client
func RegisterEagerRuleValidatingWebhook(mgr ctrl.Manager) {
	hookServer := mgr.GetWebhookServer()
	hookServer.Register("/validate-gingersnap-project-io-v1alpha1-eagercacherule", &webhook.Admission{
		Handler: &eagerRuleValidator{},
	})
}

type eagerRuleValidator struct {
	client  runtimeClient.Client
	decoder *admission.Decoder
}

var _ inject.Client = &eagerRuleValidator{}
var _ admission.Handler = &eagerRuleValidator{}

// InjectClient injects the client.
func (rv *eagerRuleValidator) InjectClient(c runtimeClient.Client) error {
	rv.client = c
	return nil
}

// InjectDecoder injects the decoder.
func (rv *eagerRuleValidator) InjectDecoder(d *admission.Decoder) error {
	rv.decoder = d
	return nil
}

func (rv *eagerRuleValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	rule := &EagerCacheRule{}
	if req.Operation == admissionv1.Create {
		err := rv.decoder.Decode(req, rule)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		err = rv.create(ctx, rule)
		if err != nil {
			var apiStatus apierrors.APIStatus
			if errors.As(err, &apiStatus) {
				return validationResponseFromStatus(false, apiStatus.Status())
			}
			return admission.Denied(err.Error())
		}
	}

	if req.Operation == admissionv1.Update {
		oldrule := &EagerCacheRule{}

		err := rv.decoder.DecodeRaw(req.Object, rule)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		err = rv.decoder.DecodeRaw(req.OldObject, oldrule)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		err = rv.update(rule, oldrule)
		if err != nil {
			var apiStatus apierrors.APIStatus
			if errors.As(err, &apiStatus) {
				return validationResponseFromStatus(false, apiStatus.Status())
			}
			return admission.Denied(err.Error())
		}
	}
	return admission.Allowed("")
}

func (rv *eagerRuleValidator) create(ctx context.Context, r *EagerCacheRule) error {
	var allErrs field.ErrorList

	spec := field.NewPath("spec")
	if r.Spec.CacheRef == nil {
		allErrs = append(allErrs, field.Required(spec.Child("cacheRef"), "'cacheRef' field must be defined"))
	} else {
		cacheRefField := spec.Child("cacheRef")
		RequireField(&allErrs, "name", r.Spec.CacheRef.Name, cacheRefField)
		RequireField(&allErrs, "namespace", r.Spec.CacheRef.Namespace, cacheRefField)
	}

	RequireField(&allErrs, "tableName", r.Spec.TableName, spec)

	if r.Spec.Key == nil {
		allErrs = append(allErrs, field.Required(spec.Child("key"), FieldMustBeDefined("key")))
	} else {
		RequireNonEmptyArray(&allErrs, "keyColumns", r.Spec.Key.KeyColumns, spec.Child("key"))
	}

	// Ensure that a EagerCacheRule CR with this cacheRef does not already exist in the cluster
	if len(allErrs) == 0 {
		cache := &CacheService{
			Name:      r.Spec.CacheRef.Name,
			Namespace: r.Spec.CacheRef.Namespace,
		}

		list := &EagerCacheRuleList{}
		listOpts := &runtimeClient.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				cache.LabelSelector(),
			),
		}
		if err := rv.client.List(ctx, list, listOpts); err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
		} else {
			for i := range list.Items {
				existingRule := &list.Items[i]
				if r.Name == existingRule.Name {
					msg := fmt.Sprintf("EagerCacheRule CR already exists for Cache '%s' with name '%s' in namespace '%s'", cache, existingRule.Name, existingRule.Namespace)
					allErrs = append(allErrs, field.Duplicate(field.NewPath("spec").Child("cacheRef"), msg))
				}
			}
		}
	}
	return StatusError(allErrs, r.Name, KindEagerCacheRule)
}

func (rv *eagerRuleValidator) update(new, old *EagerCacheRule) error {
	var allErrs field.ErrorList
	if err := EnsureRuleImmutability(&allErrs, KindEagerCacheRule, new, old); err != nil {
		return fmt.Errorf("unable to compare updated rule with existing rule: %w", err)
	}
	return StatusError(allErrs, new.Name, KindEagerCacheRule)
}
