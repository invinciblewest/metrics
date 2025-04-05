package agent

import (
	"github.com/golang/mock/gomock"
	"github.com/invinciblewest/metrics/internal/agent/collectors"
	"github.com/invinciblewest/metrics/internal/agent/collectors/mocks"
	"github.com/invinciblewest/metrics/internal/agent/senders"
	mocks2 "github.com/invinciblewest/metrics/internal/agent/senders/mocks"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	st := memstorage.NewMemStorage("", false)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mCollector := mocks.NewMockCollector(ctrl)
	mCollector.EXPECT().Collect().DoAndReturn(func() error {
		delta := int64(1)
		return st.UpdateCounter(models.Metric{
			ID:    "PollCount",
			MType: models.TypeCounter,
			Delta: &delta,
		})
	}).AnyTimes()

	mSender := mocks2.NewMockSender(ctrl)
	mSender.EXPECT().SendMetric(gomock.Any()).Return(nil).AnyTimes()

	collectorsList := []collectors.Collector{
		mCollector,
	}
	sendersList := []senders.Sender{
		mSender,
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
