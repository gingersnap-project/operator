package v1alpha1

import "fmt"

func (c *Cache) CacheService() CacheService {
	return CacheService{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

func (c *Cache) DBSyncerName() string {
	return fmt.Sprintf("%s-db-syncer", c.Name)
}
