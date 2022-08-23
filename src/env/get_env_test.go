package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvVar(t *testing.T) {
	t.Setenv("test", "1")
	val := getEnvVar("test", "test1")

	assert.Equal(t, val, "1")
}

func TestGetEnvVarDefault(t *testing.T) {
	val := getEnvVar("test", "test1")

	assert.Equal(t, val, "test1")
}
