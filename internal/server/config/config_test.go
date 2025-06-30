package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
