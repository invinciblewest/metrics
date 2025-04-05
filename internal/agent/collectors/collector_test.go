package collectors

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/invinciblewest/metrics/internal/agent/collectors/mocks"
	"github.com/invinciblewest/metrics/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	ctx := t.Context()
	st := memstorage.NewMemStorage("", false)
	t.Run("success", func(t *testing.T) {
		c := NewRuntimeCollector(st)

		err := CollectMetrics(ctx, c)
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCollector := mocks.NewMockCollector(ctrl)
		mockCollector.EXPECT().Collect(ctx).Return(errors.New("123"))

		err := CollectMetrics(ctx, mockCollector)
		assert.Error(t, err)
	})

}
