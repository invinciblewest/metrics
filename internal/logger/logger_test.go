package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		err := Initialize("test")
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		logLevel := "debug"
		err := Initialize(logLevel)
		assert.NoError(t, err)
	})
}
