package v1alpha1

import "fmt"

const (
	LabelCache          = Group + "/cache"
	LabelCacheNamespace = LabelCache + "-namespace"
)

// CacheService defines the location of the Cache resource that this LazyCacheRule should be applied to
type CacheService struct {
	// Name is the name of the Cache resource that the LazyCacheRule will be applied to
	Name string `json:"name"`
	// Namespace is the namespace in which the Cache CR belongs
	Namespace string `json:"namespace"`
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

func (s CacheService) EagerCacheConfigMap() string {
	return fmt.Sprintf("%s-eager-cm", s.Name)
}

func (s CacheService) LazyCacheConfigMap() string {
	return fmt.Sprintf("%s-lazy-cm", s.Name)
}

func (s CacheService) String() string {
	return s.Namespace + "/" + s.Name
}
