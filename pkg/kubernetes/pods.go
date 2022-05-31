package kubernetes

import corev1 "k8s.io/api/core/v1"

func GetContainer(name string, spec *corev1.PodSpec) *corev1.Container {
	for i, c := range spec.Containers {
		if c.Name == name {
			return &spec.Containers[i]
		}
	}
	return nil
}
