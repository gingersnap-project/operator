package v1alpha1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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

// RegisterLazyRuleValidatingWebhook explicitly adds the validating webhook to the Webhook Server
// This is necessary as we need to implement admission.Handler interface directly so that the request context can be
// used by the runtime client
func RegisterLazyRuleValidatingWebhook(mgr ctrl.Manager) {
	hookServer := mgr.GetWebhookServer()
	hookServer.Register("/validate-gingersnap-project-io-v1alpha1-lazycacherule", &webhook.Admission{
		Handler: &lazyRuleValidator{},
	})
}

type lazyRuleValidator struct {
	client  runtimeClient.Client
	decoder *admission.Decoder
}

var _ inject.Client = &lazyRuleValidator{}
var _ admission.Handler = &lazyRuleValidator{}

// InjectClient injects the client.
func (rv *lazyRuleValidator) InjectClient(c runtimeClient.Client) error {
	rv.client = c
	return nil
}

// InjectDecoder injects the decoder.
func (rv *lazyRuleValidator) InjectDecoder(d *admission.Decoder) error {
	rv.decoder = d
	return nil
}

func (rv *lazyRuleValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	rule := &LazyCacheRule{}
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
		oldrule := &LazyCacheRule{}

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

func (rv *lazyRuleValidator) create(ctx context.Context, r *LazyCacheRule) error {
	var allErrs field.ErrorList
	if r.Spec.CacheRef == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("cacheRef"), "'cacheRef' field must be defined"))
	} else {
		cacheRefField := field.NewPath("spec").Child("cacheRef")
		RequireField(&allErrs, "name", r.Spec.CacheRef.Name, cacheRefField)
		RequireField(&allErrs, "namespace", r.Spec.CacheRef.Namespace, cacheRefField)
	}
	RequireField(&allErrs, "query", r.Spec.Query, field.NewPath("spec"))

	// Ensure that a LazyCacheRule CR with this cacheRef does not already exist in the cluster
	if len(allErrs) == 0 {
		cache := &CacheService{
			Name:      r.Spec.CacheRef.Name,
			Namespace: r.Spec.CacheRef.Namespace,
		}

		list := &LazyCacheRuleList{}
		listOpts := &runtimeClient.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				cache.LabelSelector(),
			),
		}
		if err := rv.client.List(ctx, list, listOpts); err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
		} else if len(list.Items) > 0 {
			rule := &list.Items[0]
			msg := fmt.Sprintf("LazyCacheRule CR already exists for Cache '%s' with name '%s' in namespace '%s'", cache, rule.Name, rule.Namespace)
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec").Child("cacheRef"), msg))
		}
	}
	return StatusError(allErrs, r.Name, KindLazyCacheRule)
}

func (rv *lazyRuleValidator) update(new, old *LazyCacheRule) error {
	var allErrs field.ErrorList
	if err := EnsureRuleImmutability(&allErrs, KindLazyCacheRule, new, old); err != nil {
		return fmt.Errorf("unable to compare updated rule with existing rule: %w", err)
	}
	return StatusError(allErrs, new.Name, KindLazyCacheRule)
}

// validationResponseFromStatus returns a response for admitting a request with provided Status object.
func validationResponseFromStatus(allowed bool, status metav1.Status) admission.Response {
	return admission.Response{
		AdmissionResponse: admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result:  &status,
		},
	}
}
