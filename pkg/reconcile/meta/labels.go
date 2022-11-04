package meta

const (
	ComponentCache    = "cache"
	ComponentDBSyncer = "db-syncer"
)

func GingersnapLabels(name, component, instance string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/component":  component,
		"app.kubernetes.io/instance":   instance,
		"app.kubernetes.io/managed-by": "controller-manager",
		"app.kubernetes.io/created-by": "controller-manager",
		"app.kubernetes.io/part-of":    "gingersnap",
	}
}
