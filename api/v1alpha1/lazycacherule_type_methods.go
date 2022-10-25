package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
)

func (r *LazyCacheRule) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
}

func (r *LazyCacheRule) Filename() string {
	return fmt.Sprintf("%s_%s", r.Namespace, r.Name)
}
