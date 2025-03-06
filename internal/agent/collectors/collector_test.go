package collectors

import (
	"errors"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCollector struct{}

func (c *TestCollector) Collect() error {
	return errors.New("123")
}

func TestCollectMetrics(t *testing.T) {
	st := storage.NewMemStorage()
	t.Run("success", func(t *testing.T) {
		c := NewRuntimeCollector(st)

		err := CollectMetrics(c)
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		err := CollectMetrics(&TestCollector{})
		assert.Error(t, err)
	})

}
