package env_test

import (
	"testing"

	. "submission/env"

	"github.com/stretchr/testify/assert"
)

func TestEnvServerUrl(t *testing.T) {
	t.Setenv("SERVER_HOST", "1")
	url := GetServerHost()
	assert.Equal(t, "1", url)
}

func TestEnvServerUrlDefault(t *testing.T) {
	url := GetServerHost()
	assert.Equal(t, "http://localhost:8080", url)
}
