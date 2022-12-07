package v1alpha1

import (
	"bytes"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type CacheRule interface {
	MarshallSpec() ([]byte, error)
}

func EnsureRuleImmutability(allErrs *field.ErrorList, ruleKind string, a, b CacheRule) error {
	// We have to serialize the Proto message as proto.Equals will return false if a field which can't be compared is encountered
	aSerialized, err := a.MarshallSpec()
	if err != nil {
		return fmt.Errorf("unable to serialize first rule: %w", err)
	}

	bSerialized, err := b.MarshallSpec()
	if err != nil {
		return fmt.Errorf("unable to serialize second rule: %w", err)
	}

	if !bytes.Equal(aSerialized, bSerialized) {
		detail := fmt.Sprintf("The %s spec is immutable and cannot be updated after initial creation", ruleKind)
		*allErrs = append(*allErrs, field.Forbidden(field.NewPath("spec"), detail))
	}
	return nil
}

func StatusError(allErrs field.ErrorList, name, kind string) error {
	if len(allErrs) != 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: GroupVersion.Group, Kind: kind}, name, allErrs,
		)
	}
	return nil
}

func RequireField(allErrs *field.ErrorList, name, value string, p *field.Path) {
	if value == "" {
		*allErrs = append(*allErrs, field.Required(p.Child(name), emptyFieldDetail(name)))
	}
}

func RequireNonEmptyArray(allErrs *field.ErrorList, name string, value []string, p *field.Path) {
	if len(value) == 0 {
		*allErrs = append(*allErrs, field.Required(p.Child(name), emptyFieldDetail(name)))
	}
}

func FieldMustBeDefined(field string) string {
	return fmt.Sprintf("'%s' field must be defined", field)
}

func emptyFieldDetail(field string) string {
	return fmt.Sprintf("'%s' field must not be empty", field)
}
