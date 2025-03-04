package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
