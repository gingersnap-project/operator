package v1alpha1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (c *Cache) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gingersnap-project-io-v1alpha1-cache,mutating=true,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=caches,verbs=create;update,versions=v1alpha1,name=mcache.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Cache{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (c *Cache) Default() {
	if c.Spec.Deployment == nil {
		c.Spec.Deployment = &CacheDeploymentSpec{
			Type: CacheDeploymentType_LOCAL,
		}
	}

	if c.Cluster() {
		if c.Spec.Deployment.Replicas < 1 {
			c.Spec.Deployment.Replicas = 1
		}
	}
}

//+kubebuilder:webhook:path=/validate-gingersnap-project-io-v1alpha1-cache,mutating=false,failurePolicy=fail,sideEffects=None,groups=gingersnap-project.io,resources=caches,verbs=create;update,versions=v1alpha1,name=vcache.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Cache{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (c *Cache) ValidateCreate() error {
	return c.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (c *Cache) ValidateUpdate(_ runtime.Object) error {
	if err := c.validate(); err != nil {
		return err
	}

	var allErrs field.ErrorList
	return c.statusError(allErrs)
}

func (c *Cache) validate() error {
	var allErrs field.ErrorList

	validateResources(&allErrs, field.NewPath("spec").Child("deployment").Child("resources"), c.Spec.Deployment.Resources)

	if c.Spec.DbSyncer != nil {
		validateResources(&allErrs, field.NewPath("spec").Child("dbSyncer").Child("resources"), c.Spec.DbSyncer.Resources)
	}

	ds := c.Spec.DataSource
	if ds == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("dataSource"), "A dataSource must be defined"))
	} else {
		if ds.DbType == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("dataSource").Child("dbType"), "A dataSource dbType must be defined"))
		}

		if ds.SecretRef != nil && ds.ServiceProviderRef != nil {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec").Child("dataSource"), "At most one of ['secretRef', 'serviceProviderRef'] must be configured"))
		} else if ds.SecretRef == nil && ds.ServiceProviderRef == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("dataSource"), "'secretRef' OR 'serviceProviderRef' must be supplied"))
		} else {
			if ds.SecretRef != nil {
				requireField(&allErrs, "name", ds.SecretRef.Name, field.NewPath("spec").Child("dataSource").Child("secretRef"))
			} else if ds.ServiceProviderRef != nil {
				root := field.NewPath("spec").Child("dataSource").Child("serviceProviderRef")
				requireField(&allErrs, "apiVersion", ds.ServiceProviderRef.ApiVersion, root)
				requireField(&allErrs, "kind", ds.ServiceProviderRef.Kind, root)
				requireField(&allErrs, "name", ds.ServiceProviderRef.Name, root)
			}
		}
	}
	return c.statusError(allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (c *Cache) ValidateDelete() error {
	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (c *Cache) statusError(allErrs field.ErrorList) error {
	if len(allErrs) != 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: GroupVersion.Group, Kind: KindCache},
			c.Name, allErrs)
	}
	return nil
}

func requireField(allErrs *field.ErrorList, name, value string, p *field.Path) {
	if value == "" {
		*allErrs = append(*allErrs, field.Required(p.Child(name), fmt.Sprintf("'%s' field must not be empty", name)))
	}
}

func validateResources(allErrs *field.ErrorList, p *field.Path, r *Resources) {
	if r == nil {
		return
	}

	parse := func(quantity, resourceType, name string) {
		_, err := resource.ParseQuantity(quantity)
		if err != nil {
			*allErrs = append(*allErrs, field.Invalid(p.Child(resourceType).Child(name), quantity, err.Error()))
		}
	}

	if r.Requests != nil {
		parse(r.Requests.Cpu, "requests", "cpu")
		parse(r.Requests.Memory, "requests", "memory")
	}

	if r.Limits != nil {
		parse(r.Limits.Cpu, "limits", "cpu")
		parse(r.Limits.Memory, "limits", "memory")
	}
}
