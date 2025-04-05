package collectors

import (
	"context"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"runtime"
)

type RuntimeCollector struct {
	st storage.Storage
}

func NewRuntimeCollector(st storage.Storage) *RuntimeCollector {
	return &RuntimeCollector{
		st: st,
	}
}

func (c *RuntimeCollector) Collect(ctx context.Context) error {
	logger.Log.Info("collecting metrics...",
		zap.String("collector", "runtime"),
	)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricsIds := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}

	value := reflect.ValueOf(memStats)
	for _, key := range metricsIds {
		field := value.FieldByName(key)
		if field.IsValid() {
			var fieldValue float64
			if field.CanFloat() {
				fieldValue = field.Float()
			} else if field.CanUint() {
				fieldValue = float64(field.Uint())
			}
			err := c.st.UpdateGauge(ctx, models.Metric{
				ID:    key,
				MType: models.TypeGauge,
				Value: &fieldValue,
			})
			if err != nil {
				return err
			}
		}
	}

	randomFloat := rand.Float64()
	err := c.st.UpdateGauge(ctx, models.Metric{
		ID:    "RandomValue",
		MType: models.TypeGauge,
		Value: &randomFloat,
	})
	if err != nil {
		return err
	}
	counterValue := int64(1)
	err = c.st.UpdateCounter(ctx, models.Metric{
		ID:    "PollCount",
		MType: models.TypeCounter,
		Delta: &counterValue,
	})
	if err != nil {
		return err
	}

	pc, err := c.st.GetCounter(ctx, "PollCount")
	if err != nil {
		return err
	}
	logger.Log.Info("poll is collected",
		zap.String("collector", "runtime"),
		zap.Int64("poll", *pc.Delta),
	)
	return nil
}
