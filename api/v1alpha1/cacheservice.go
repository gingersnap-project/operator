package v1alpha1

const (
	LabelCache          = Group + "/cache"
	LabelCacheNamespace = LabelCache + "-namespace"
)

func CacheServiceLabelsExist(m map[string]string) bool {
	_, cacheExists := m[LabelCache]
	_, namespaceExists := m[LabelCacheNamespace]
	return cacheExists && namespaceExists
}

func CacheServiceFromLabels(m map[string]string) CacheService {
	return CacheService{
		Namespace: m[LabelCacheNamespace],
		Name:      m[LabelCache],
	}
}

func (s CacheService) ApplyLabels(m map[string]string) {
	m[LabelCache] = s.Name
	m[LabelCacheNamespace] = s.Namespace
}

func (s CacheService) LabelSelector() map[string]string {
	return map[string]string{
		LabelCache:          s.Name,
		LabelCacheNamespace: s.Namespace,
	}
}

func (s CacheService) String() string {
	return s.Namespace + "/" + s.Name
}
