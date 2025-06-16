package collectors

import (
	"context"
	"fmt"

	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

// GopsutilCollector собирает метрики с помощью библиотеки gopsutil.
type GopsutilCollector struct {
	st storage.Storage
}

// NewGopsutilCollector создает новый экземпляр GopsutilCollector с заданным хранилищем.
func NewGopsutilCollector(st storage.Storage) *GopsutilCollector {
	return &GopsutilCollector{
		st: st,
	}
}

// Collect собирает метрики о памяти и загрузке процессора и сохраняет их в хранилище.
func (c *GopsutilCollector) Collect(ctx context.Context) error {
	logger.Log.Info("collecting metrics...",
		zap.String("collector", "gopsutil"),
	)

	virtualMemoryStat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	totalMemory := float64(virtualMemoryStat.Total)
	freeMemory := float64(virtualMemoryStat.Free)

	metrics := []models.Metric{
		{
			ID:    "TotalMemory",
			MType: models.TypeGauge,
			Value: &totalMemory,
		},
		{
			ID:    "FreeMemory",
			MType: models.TypeGauge,
			Value: &freeMemory,
		},
	}

	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}

	for i, percent := range cpuPercent {
		metrics = append(metrics, models.Metric{
			ID:    fmt.Sprintf("CPUutilization%d", i),
			MType: models.TypeGauge,
			Value: &percent,
		})
	}

	err = c.st.UpdateBatch(ctx, metrics)
	if err != nil {
		return err
	}

	logger.Log.Info("poll is collected",
		zap.String("collector", "gopsutil"),
	)
	return nil
}
