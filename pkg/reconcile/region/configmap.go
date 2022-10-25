package region

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

func configMapLabels(cacheService v1alpha1.CacheService) map[string]string {
	labels := map[string]string{
		"app.kubernetes.io/name":       "gingersnap",
		"app.kubernetes.io/managed-by": "controller-manager",
		"app.kubernetes.io/created-by": "controller-manager",
	}
	cacheService.ApplyLabels(labels)
	return labels
}

func UpdateConfigMaps(r *v1alpha1.CacheRegion, ctx reconcile.Context) {
	cmList := &corev1.ConfigMapList{}
	labels := configMapLabels(r.Spec.Cache)
	if err := ctx.Client().List(labels, cmList); err != nil {
		ctx.Requeue(fmt.Errorf("unable to list existing ConfigMaps: %w", err))
		return
	}

	bytes, err := marshallRegion(r)
	if err != nil {
		ctx.Requeue(err)
		return
	}

	key := r.Filename()
	for _, cm := range cmList.Items {
		if cm.BinaryData == nil {
			cm.BinaryData = map[string][]byte{}
		}
		cm.BinaryData[key] = bytes
		if err := ctx.Client().WithNamespace(cm.Namespace).Update(&cm); err != nil {
			ctx.Requeue(fmt.Errorf("unable to update ConfigMap '%s:%s': %v", cm.Namespace, cm.Name, err))
			return
		}
	}
}

func marshallRegion(region *v1alpha1.CacheRegion) ([]byte, error) {
	bytes, err := yaml.Marshal(region)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall CacheRegion '%s", region.Filename())
	}
	return bytes, nil
}
