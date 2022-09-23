package region

import (
	"fmt"

	"github.com/gingersnap-project/operator/api/v1alpha1"
	"github.com/gingersnap-project/operator/pkg/kubernetes/client"
	"github.com/gingersnap-project/operator/pkg/reconcile"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func CreateConfigMap(namespace string, service v1alpha1.CacheService, k8sClient client.Client) (*corev1.ConfigMap, error) {
	k8sClient = k8sClient.WithNamespace(namespace)
	labels := configMapLabels(service)
	configMapData, err := ConfigMapData(service, k8sClient)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize CacheRegion ConfigMap data: %w", err)
	}

	// Create ConfigMap using corev1 type so that we make use of GenerateName
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: "gingersnap-regions-",
			Labels:       labels,
		},
		BinaryData: configMapData,
	}

	if err := k8sClient.Create(configMap); err != nil && !errors.IsAlreadyExists(err) {
		return nil, fmt.Errorf("unable to create CacheRegion ConfigMap: %w", err)
	}

	// We must loop until the created ConfigMap can be retrieved so that the generated ConfigMap name can be used by
	// the sidecar container.
	// Any error other than NotFound is treated as a failure.
	// If for some unforeseen reason the ConfigMap is never found, the Webhook will timeout and the context will be cancelled.
	for {
		if configMap, err = FindExistingConfigMap(namespace, service, k8sClient); err != nil {
			return nil, err
		} else if configMap != nil {
			return configMap, err
		}
	}
}

func FindExistingConfigMap(namespace string, service v1alpha1.CacheService, client client.Client) (*corev1.ConfigMap, error) {
	client = client.WithNamespace(namespace)
	labels := configMapLabels(service)
	cmList := &corev1.ConfigMapList{}
	if err := client.List(labels, cmList); err != nil {
		return nil, fmt.Errorf("unable to list existing ConfigMaps: %w", err)
	}

	switch numConfigMaps := len(cmList.Items); numConfigMaps {
	case 0:
		return nil, nil
	case 1:
		return &cmList.Items[0], nil
	default:
		var configMapNames []string
		for i, cm := range cmList.Items {
			configMapNames[i] = cm.Name
		}
		return nil, fmt.Errorf("expected one ConfigMap, found %d: %s", numConfigMaps, configMapNames)
	}
}

// ConfigMapData retrieves all CacheRegion resources related to the passed CacheService, iterates over their Region
// definitions and adds their json representation to a string:[]byte map that can be used by a ConfigMap
func ConfigMapData(service v1alpha1.CacheService, k8sClient client.Client) (map[string][]byte, error) {
	regionList := &v1alpha1.CacheRegionList{}
	if err := k8sClient.List(service.LabelSelector(), regionList, client.ClusterScoped); err != nil {
		return nil, fmt.Errorf("unable to list CacheRegion: %w", err)
	}

	regionMap := make(map[string][]byte, len(regionList.Items))
	for _, region := range regionList.Items {
		if bytes, err := marshallRegion(&region); err != nil {
			return nil, err
		} else {
			regionMap[region.Filename()] = bytes
		}
	}
	return regionMap, nil
}

func marshallRegion(region *v1alpha1.CacheRegion) ([]byte, error) {
	bytes, err := yaml.Marshal(region)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall CacheRegion '%s", region.Filename())
	}
	return bytes, nil
}
