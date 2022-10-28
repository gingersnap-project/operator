package lazy

import (
	"encoding/json"
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var finalizer = schema.GroupKind{Group: v1alpha1.Group, Kind: v1alpha1.KindLazyCacheRule}.String()

func LoadCache(r *v1alpha1.LazyCacheRule, ctx *Context) {
	ctx.Cache = &v1alpha1.Cache{}
	cacheRef := r.Spec.Cache
	err := ctx.Client().
		WithNamespace(cacheRef.Namespace).
		Load(cacheRef.Name, ctx.Cache)

	if err != nil {
		// TODO set status !Ready condition
		ctx.Requeue(fmt.Errorf("unable to load Cache CR '%s': %w", cacheRef, err))
	}
}

func AddRuleToConfigMap(r *v1alpha1.LazyCacheRule, ctx *Context) {
	existingConfigMap, err := loadRuleConfigMap(ctx.Cache.CacheService(), ctx)
	if err != nil {
		ctx.Requeue(err)
		return
	}

	var data map[string]string
	if existingConfigMap == nil {
		data = make(map[string]string, 1)
	} else {
		data = existingConfigMap.Data
	}

	bytes, err := json.Marshal(r)
	if err != nil {
		ctx.Requeue(fmt.Errorf("unable to marshall rule: %w", err))
		return
	}
	data[r.Filename()] = string(bytes[:])

	applyRuleConfigMap(data, r, ctx)
}

func AddFinalizer(r *v1alpha1.LazyCacheRule, ctx *Context) {
	if controllerutil.AddFinalizer(r, finalizer) {
		if err := ctx.Client().Update(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to add finalizer: %w", err))
		}
	}
}

func loadRuleConfigMap(cache v1alpha1.CacheService, ctx *Context) (*apicorev1.ConfigMap, error) {
	existingConfigMap := &apicorev1.ConfigMap{}
	cmName := cache.LazyCacheConfigMap()

	err := ctx.Client().
		WithNamespace(cache.Namespace).
		Load(cmName, existingConfigMap)

	if runtimeClient.IgnoreNotFound(err) != nil {
		return nil, fmt.Errorf("unable to load ConfigMap '%s': %w", cmName, err)
	}

	if errors.IsNotFound(err) {
		return nil, nil
	}
	return existingConfigMap, nil
}

func applyRuleConfigMap(data map[string]string, r *v1alpha1.LazyCacheRule, ctx *Context) {
	cache := ctx.Cache
	cmName := cache.CacheService().LazyCacheConfigMap()

	labels := configMapLabels(r.Spec.Cache)
	cm := corev1.
		ConfigMap(cmName, cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(client.OwnerReference(cache)).
		WithData(data)

	if err := ctx.Client().Apply(cm); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply '%s' ConfigMap: %w", cmName, err))
	}
}

func configMapLabels(cacheService v1alpha1.CacheService) map[string]string {
	labels := map[string]string{
		"app.kubernetes.io/name":       "gingersnap",
		"app.kubernetes.io/managed-by": "controller-manager",
		"app.kubernetes.io/created-by": "controller-manager",
	}
	cacheService.ApplyLabels(labels)
	return labels
}
