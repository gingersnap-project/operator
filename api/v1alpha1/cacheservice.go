package v1alpha1

const (
	LabelCache          = Group + "/cache"
	LabelCacheNamespace = LabelCache + "-namespace"
)

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
