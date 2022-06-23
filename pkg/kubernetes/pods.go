package kubernetes

import corev1 "k8s.io/api/core/v1"

func Container(name string, spec *corev1.PodSpec) *corev1.Container {
	for i, c := range spec.Containers {
		if c.Name == name {
			return &spec.Containers[i]
		}
	}
	return nil
}

func Volume(name string, spec *corev1.PodSpec) *corev1.Volume {
	for i, v := range spec.Volumes {
		if v.Name == name {
			return &spec.Volumes[i]
		}
	}
	return nil
}
