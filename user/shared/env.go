package shared

import (
	"os"
)

// GetStringEnvVar gets the env var as a string
func GetStringEnvVar(varName string, defaultValue string) string {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		return defaultValue
	}

	return val
}
