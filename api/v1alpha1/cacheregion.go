package v1alpha1

import "strings"

const (
	AnnotationRegions = Group + "/regions"
)

func CacheRegionsFromAnnotations(m map[string]string) []string {
	annotation, ok := m[AnnotationRegions]
	if !ok {
		return []string{}
	}
	return strings.Split(annotation, ",")
}
