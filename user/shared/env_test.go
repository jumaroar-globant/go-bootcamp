package shared

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStringEnvVarDefaultValue(t *testing.T) {
	c := require.New(t)

	defaultValue := "val"
	c.Equal(defaultValue, GetStringEnvVar("GET_STRING", defaultValue))
}

func TestGetStringEnvVarCustomValue(t *testing.T) {
	c := require.New(t)

	withTestEnv("custom", func(varName string) {
		c.Equal("custom", GetStringEnvVar(varName, ""))
	})
}

func withTestEnv(val interface{}, cb func(varName string)) {
	varName := fmt.Sprintf("TEST_%d_%d", rand.Intn(math.MaxInt32), rand.Intn(math.MaxInt32))
	_ = os.Setenv(varName, fmt.Sprintf("%v", val))
	defer func() {
		_ = os.Unsetenv(varName)
	}()

	cb(varName)
}
