package v1alpha1

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (c *Cache) CacheService() CacheService {
	return CacheService{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

func (c *Cache) DBSyncerName() string {
	return fmt.Sprintf("%s-db-syncer", c.Name)
}

func (c *Cache) DeploymentLimits() v1.ResourceList {
	if c.Spec.Deployment != nil && c.Spec.Deployment.Resources != nil && c.Spec.Deployment.Resources.Limits != nil {
		return resourceList(c.Spec.Deployment.Resources.Limits)
	}
	return nil
}

func (c *Cache) DeploymentRequests() v1.ResourceList {
	if c.Spec.Deployment != nil && c.Spec.Deployment.Resources != nil && c.Spec.Deployment.Resources.Requests != nil {
		return resourceList(c.Spec.Deployment.Resources.Requests)
	}
	return nil
}

func (c *Cache) DBSyncerLimits() v1.ResourceList {
	if c.Spec.DbSyncer != nil && c.Spec.DbSyncer.Resources != nil && c.Spec.DbSyncer.Resources.Limits != nil {
		return resourceList(c.Spec.DbSyncer.Resources.Limits)
	}
	return nil
}

func (c *Cache) DBSyncerRequests() v1.ResourceList {
	if c.Spec.DbSyncer != nil && c.Spec.DbSyncer.Resources != nil && c.Spec.DbSyncer.Resources.Requests != nil {
		return resourceList(c.Spec.DbSyncer.Resources.Requests)
	}
	return nil
}

func resourceList(rq *ResourceQuantity) v1.ResourceList {
	// MustParse should never throw a panic as the webhook has already verified that the quantity is valid
	return map[v1.ResourceName]resource.Quantity{
		v1.ResourceCPU:    resource.MustParse(rq.Cpu),
		v1.ResourceMemory: resource.MustParse(rq.Memory),
	}
}
