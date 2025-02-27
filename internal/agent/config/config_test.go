package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()
	assert.NotNil(t, cfg)
}
