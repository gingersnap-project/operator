package e2e

import (
	"os"
)

const (
	DefaultNamespace = "gingersnap-operator-system"
)

var (
	OperatorNamespace = EnvWithDefault("TEST_OPERATOR_NAMESPACE", DefaultNamespace)
	Namespace         = EnvWithDefault("TEST_NAMESPACE", DefaultNamespace)
	OutputDir         = EnvWithDefault("TEST_OUTPUT_DIR", os.TempDir()+"/gingersnap-operator")

	MultiNamespace = Namespace != OperatorNamespace
)

// WithDefault return value if not empty else return defValue
func WithDefault(value, defValue string) string {
	if value == "" {
		return defValue
	}
	return value
}

// EnvWithDefault return os.Getenv(name) if exists else return defValue
func EnvWithDefault(name, defValue string) string {
	return WithDefault(os.Getenv(name), defValue)
}
