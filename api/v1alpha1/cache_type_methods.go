package v1alpha1

import "fmt"

func (c *Cache) CacheService() CacheService {
	return CacheService{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

func (c *Cache) ConfigurationSecret() string {
	return c.Name
}

func (c *Cache) SvcName() string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", c.Name, c.Namespace)
}
