package agent

import (
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockCollector struct {
	st storage.Storage
}

func (c *MockCollector) Collect() error {
	delta := int64(1)
	c.st.UpdateCounter(models.Metrics{
		ID:    "PollCount",
		MType: models.TypeCounter,
		Delta: &delta,
	})
	return nil
}

type MockSender struct{}

func (s *MockSender) Send(metrics models.Metrics) error {
	return nil
}

func TestNewAgent(t *testing.T) {
	st := storage.NewMemStorage("", false)
	collectorsList := []collectors.Collector{
		&MockCollector{
			st: st,
		},
	}
	sendersList := []senders.Sender{
		&MockSender{},
	}

	go func() {
		agent := NewAgent(st, collectorsList, sendersList, 1, 2)
		err := agent.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(2 * time.Second)

	_, err := st.GetCounter("PollCount")
	assert.NoError(t, err)
}
