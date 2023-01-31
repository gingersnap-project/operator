package v1alpha1

import (
	"fmt"

	"github.com/gingersnap-project/operator/pkg/images"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Cache) CacheManagerImage() string {
	switch *c.Spec.DataSource.DbType {
	case DBType_MYSQL_8:
		return images.CacheManagerMySQL
	case DBType_POSTGRES_14:
		return images.CacheManagerPostgres
	case DBType_SQL_SERVER_2019:
		return images.CacheManagerMSSQL
	}
	return ""
}

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

func (c *Cache) Local() bool {
	return c.Spec.Deployment.Type == CacheDeploymentType_LOCAL
}

func (c *Cache) Cluster() bool {
	return c.Spec.Deployment.Type == CacheDeploymentType_CLUSTER
}

func (c *Cache) Condition(condition CacheConditionType) CacheCondition {
	for _, existing := range c.Status.Conditions {
		if existing.Type == condition {
			return existing
		}
	}
	// Absence of condition means `False` value
	return CacheCondition{Type: condition, Status: metav1.ConditionFalse}
}

func (c *Cache) SetCondition(condition CacheCondition) (updated bool) {
	for idx := range c.Status.Conditions {
		c := &c.Status.Conditions[idx]
		if c.Type == condition.Type {
			if c.Status != condition.Status {
				c.Status = condition.Status
				updated = true
			}
			if c.Message != condition.Message {
				c.Message = condition.Message
				updated = true
			}
			return updated
		}
	}
	c.Status.Conditions = append(c.Status.Conditions, CacheCondition{
		Type:    condition.Type,
		Status:  condition.Status,
		Message: condition.Message,
	})
	return true
}

func (x CacheDeploymentType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", CacheDeploymentType_name[int32(x)])), nil
}

func (x *CacheDeploymentType) UnmarshalJSON(b []byte) error {
	*x = CacheDeploymentType(CacheDeploymentType_value[string(b[1:len(b)-1])])
	return nil
}

func (dbType *DBType) ServiceBinding() string {
	switch *dbType {
	case DBType_MYSQL_8:
		return "mysql"
	case DBType_POSTGRES_14:
		return "postgresql"
	case DBType_SQL_SERVER_2019:
		return "sqlserver"
	}
	return ""
}
