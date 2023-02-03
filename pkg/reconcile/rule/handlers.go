package rule

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func AddFinalizer(r CacheRule, ctx *Context) {
	if controllerutil.AddFinalizer(r, r.Finalizer()) {
		if err := ctx.Client().Update(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to add finalizer: %w", err))
		}
	}
}

func RemoveFinalizer(r CacheRule, ctx *Context) {
	if controllerutil.RemoveFinalizer(r, r.Finalizer()) {
		if err := ctx.Client().Update(r); err != nil {
			ctx.Requeue(fmt.Errorf("unable to remove finalizer: %w", err))
		}
	}
}

func ApplyRuleConfigMap(rule CacheRule, ctx *Context) {
	cache := rule.CacheService()
	existingConfigMap, err := loadConfigMap(rule.ConfigMap(), cache.Namespace, ctx)
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

	bytes, err := rule.MarshallSpec()
	if err != nil {
		ctx.Requeue(fmt.Errorf("unable to marshall rule: %w", err))
		return
	}
	data[rule.GetName()] = string(bytes[:])

	labels := configMapLabels(cache)
	cm := corev1.
		ConfigMap(rule.ConfigMap(), cache.Namespace).
		WithLabels(labels).
		WithOwnerReferences(client.OwnerReference(ctx.Cache)).
		WithData(data)

	if err := ctx.Client().Apply(cm); err != nil {
		ctx.Requeue(fmt.Errorf("unable to apply '%s' ConfigMap: %w", rule.ConfigMap(), err))
	}
}

func RemoveRuleFromConfigMap(rule CacheRule, ctx *Context) {
	cache := rule.CacheService()
	existingConfigMap, err := loadConfigMap(rule.ConfigMap(), cache.Namespace, ctx)
	if err != nil {
		ctx.Requeue(err)
		return
	}

	if existingConfigMap != nil {
		delete(existingConfigMap.Data, rule.GetName())
		if err := ctx.Client().Update(existingConfigMap); runtimeClient.IgnoreNotFound(err) != nil {
			ctx.Requeue(fmt.Errorf("unable to remove '%s' from ConfigMap: %w", rule.GetName(), err))
		}
	}
}

func configMapLabels(cacheService v1alpha1.CacheService) map[string]string {
	labels := map[string]string{
		"app.kubernetes.io/name":       "gingersnap",
		"app.kubernetes.io/managed-by": "controller-manager",
		"app.kubernetes.io/created-by": "controller-manager",
	}
	cacheService.ApplyLabelsToMap(labels)
	return labels
}

func loadConfigMap(name, namespace string, ctx reconcile.Context) (*apicorev1.ConfigMap, error) {
	existingConfigMap := &apicorev1.ConfigMap{}

	err := ctx.Client().
		WithNamespace(namespace).
		Load(name, existingConfigMap)

	if runtimeClient.IgnoreNotFound(err) != nil {
		return nil, fmt.Errorf("unable to load ConfigMap '%s': %w", name, err)
	}

	if errors.IsNotFound(err) {
		return nil, nil
	}
	return existingConfigMap, nil
}
