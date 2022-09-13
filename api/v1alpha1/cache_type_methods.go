package v1alpha1

func (c *Cache) CacheService() CacheService {
	return CacheService{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

func (c *Cache) ConfigurationSecret() string {
	return c.Name
}
