package main

import (
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockCollector struct{}

func (c *MockCollector) Collect(st *storage.MemStorage) error {
	st.UpdateCounter("PollCount", 1)
	return nil
}

type MockSender struct{}

func (s *MockSender) Send(mType string, mName string, mValue string) error {
	return nil
}

func TestRunAgent(t *testing.T) {
	st := storage.NewMemStorage()
	collectorsList := []collectors.Collector{
		&MockCollector{},
	}
	sendersList := []senders.Sender{
		&MockSender{},
	}

	go func() {
		runAgent(st, collectorsList, sendersList, 1, 2)
	}()

	time.Sleep(1 * time.Second)

	assert.Equal(t, int64(2), st.GetCounter("PollCount"))
}
