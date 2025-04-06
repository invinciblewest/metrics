package memstorage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
	"os"
	"sync"
)

type MemStorage struct {
	Gauges   storage.GaugeList   `json:"gauges"`
	Counters storage.CounterList `json:"counters"`
	path     string
	syncSave bool
	mu       sync.RWMutex
}

func NewMemStorage(path string, syncSave bool) *MemStorage {
	return &MemStorage{
		Gauges:   make(storage.GaugeList),
		Counters: make(storage.CounterList),
		path:     path,
		syncSave: syncSave,
	}
}

func (st *MemStorage) UpdateGauge(ctx context.Context, metric models.Metric) error {
	if metric.MType != models.TypeGauge {
		return storage.ErrWrongType
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	st.Gauges[metric.ID] = metric
	if st.syncSave {
		return st.Save(ctx)
	}
	return nil
}

func (st *MemStorage) GetGauge(ctx context.Context, id string) (models.Metric, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	value, exists := st.Gauges[id]
	if exists {
		return value, nil
	} else {
		return models.Metric{}, storage.ErrNotFound
	}
}

func (st *MemStorage) GetGaugeList(ctx context.Context) storage.GaugeList {
	st.mu.Lock()
	defer st.mu.Unlock()

	return st.Gauges
}

func (st *MemStorage) UpdateCounter(ctx context.Context, metric models.Metric) error {
	if metric.MType != models.TypeCounter {
		return storage.ErrWrongType
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	currentMetric, exists := st.Counters[metric.ID]
	if exists {
		*metric.Delta += *currentMetric.Delta
	}

	st.Counters[metric.ID] = metric

	if st.syncSave {
		return st.Save(ctx)
	}
	return nil
}

func (st *MemStorage) GetCounter(ctx context.Context, id string) (models.Metric, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	value, exists := st.Counters[id]
	if exists {
		return value, nil
	} else {
		return models.Metric{}, storage.ErrNotFound
	}
}

func (st *MemStorage) GetCounterList(ctx context.Context) storage.CounterList {
	st.mu.Lock()
	defer st.mu.Unlock()

	return st.Counters
}

func (st *MemStorage) UpdateBatch(ctx context.Context, metrics []models.Metric) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	for _, metric := range metrics {
		switch metric.MType {
		case models.TypeGauge:
			st.Gauges[metric.ID] = metric
		case models.TypeCounter:
			currentMetric, exists := st.Counters[metric.ID]
			if exists {
				*metric.Delta += *currentMetric.Delta
			}
			st.Counters[metric.ID] = metric
		default:
			return storage.ErrWrongType
		}
	}

	return nil
}

func (st *MemStorage) Save(ctx context.Context) error {
	if st.path == "" {
		return nil
	}
	logger.Log.Info("saving storage...", zap.String("storage", st.path))

	st.mu.Lock()
	defer st.mu.Unlock()

	file, err := os.OpenFile(st.path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer closeFile(file)

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(st); err != nil { // Передаем саму структуру
		return err
	}
	if _, err = file.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (st *MemStorage) Load(ctx context.Context) error {
	if st.path == "" {
		return nil
	}
	logger.Log.Info("loading storage...", zap.String("storage", st.path))

	_, err := os.Stat(st.path)
	if os.IsNotExist(err) {
		logger.Log.Info("storage file not exists", zap.String("path", st.path))
		return nil
	}

	file, err := os.OpenFile(st.path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer closeFile(file)

	st.mu.Lock()
	defer st.mu.Unlock()

	if err = json.NewDecoder(file).Decode(&st); err != nil {
		return err
	}

	return nil
}

func (st *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (st *MemStorage) Close(ctx context.Context) {

}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		logger.Log.Error("close file error", zap.Error(err))
	}
}
