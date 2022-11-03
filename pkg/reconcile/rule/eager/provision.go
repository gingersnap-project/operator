package eager

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile/rule"
)

func LoadCache(r *v1alpha1.EagerCacheRule, ctx *rule.Context) {
	cacheRef := r.Spec.Cache
	cache := &v1alpha1.Cache{}
	err := ctx.Client().
		WithNamespace(cacheRef.Namespace).
		Load(cacheRef.Name, cache)

	if err != nil {
		// TODO set status !Ready condition
		ctx.Requeue(fmt.Errorf("unable to load Cache CR '%s': %w", cacheRef, err))
	}
	ctx.Cache = cache
}
