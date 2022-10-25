package lazy

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func RemoveFinalizer(r *v1alpha1.LazyCacheRule, ctx *Context) {
	if controllerutil.RemoveFinalizer(r, finalizer) {
		if err := ctx.Client().Update(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to remove finalizer: %w", err))
		}
	}
}

func RemoveRuleFromConfigMap(r *v1alpha1.LazyCacheRule, ctx *Context) {
	existingConfigMap, err := loadRuleConfigMap(r.Spec.Cache, ctx)
	if err != nil {
		ctx.Requeue(err)
		return
	}

	if existingConfigMap != nil {
		delete(existingConfigMap.BinaryData, r.Filename())
		if err := ctx.Client().Update(existingConfigMap); runtimeClient.IgnoreNotFound(err) != nil {
			ctx.Requeue(fmt.Errorf("unable to remove '%s' from ConfigMap: %w", r.Filename(), err))
		}
	}
}
