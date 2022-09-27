package e2e

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DefaultNamespace = "gingersnap-operator-system"
)

var (
	OperatorNamespace    = EnvWithDefault("TEST_OPERATOR_NAMESPACE", DefaultNamespace)
	Namespace            = EnvWithDefault("TEST_NAMESPACE", DefaultNamespace)
	CleanupTestNamespace = EnvWithDefaultBool("TEST_NAMESPACE_DELETE", true)
	OutputDir            = EnvWithDefault("TEST_OUTPUT_DIR", os.TempDir()+"/gingersnap-operator")

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

// EnvWithDefaultBool return os.Getenv(name) if exists else return defValue
func EnvWithDefaultBool(name string, defValue bool) bool {
	env := os.Getenv(name)
	if env == "" {
		return defValue
	}

	value, err := strconv.ParseBool(env)
	if err != nil {
		panic(fmt.Errorf("invalid bool value for env '%s': %w", name, err))
	}
	return value
}
