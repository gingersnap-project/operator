package v1alpha1

func (c *Cache) CacheService() CacheService {
	return CacheService{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}
