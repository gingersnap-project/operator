package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (r *EagerCacheRule) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
}

func (r *EagerCacheRule) Filename() string {
	return fmt.Sprintf("%s_%s", r.Namespace, r.Name)
}

func (r *EagerCacheRule) Finalizer() string {
	return schema.GroupKind{Group: Group, Kind: KindEagerCacheRule}.String()
}

func (r *EagerCacheRule) CacheService() CacheService {
	return CacheService{
		Name:      r.Spec.CacheRef.Name,
		Namespace: r.Spec.CacheRef.Namespace,
	}
}

func (r *EagerCacheRule) ConfigMap() string {
	return r.CacheService().EagerCacheConfigMap()
}
