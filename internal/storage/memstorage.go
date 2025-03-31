package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"go.uber.org/zap"
	"os"
	"sync"
)

type MemStorage struct {
	Gauges   GaugeList   `json:"gauges"`
	Counters CounterList `json:"counters"`
	path     string
	syncSave bool
	mu       sync.RWMutex
}

func NewMemStorage(path string, syncSave bool) *MemStorage {
	return &MemStorage{
		Gauges:   make(GaugeList),
		Counters: make(CounterList),
		path:     path,
		syncSave: syncSave,
	}
}

func (st *MemStorage) UpdateGauge(metric models.Metric) error {
	if metric.MType != models.TypeGauge {
		return errors.New("wrong type")
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	st.Gauges[metric.ID] = metric
	if st.syncSave {
		return st.Save()
	}
	return nil
}

func (st *MemStorage) GetGauge(id string) (models.Metric, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	value, exists := st.Gauges[id]
	if exists {
		return value, nil
	} else {
		return models.Metric{}, ErrNotFound
	}
}

func (st *MemStorage) GetGaugeList() GaugeList {
	st.mu.Lock()
	defer st.mu.Unlock()

	return st.Gauges
}

func (st *MemStorage) UpdateCounter(metric models.Metric) error {
	if metric.MType != models.TypeCounter {
		return errors.New("wrong type")
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	currentMetric, exists := st.Counters[metric.ID]
	if exists {
		*metric.Delta += *currentMetric.Delta
	}

	st.Counters[metric.ID] = metric

	if st.syncSave {
		return st.Save()
	}
	return nil
}

func (st *MemStorage) GetCounter(id string) (models.Metric, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	value, exists := st.Counters[id]
	if exists {
		return value, nil
	} else {
		return models.Metric{}, ErrNotFound
	}
}

func (st *MemStorage) GetCounterList() CounterList {
	st.mu.Lock()
	defer st.mu.Unlock()

	return st.Counters
}

func (st *MemStorage) Save() error {
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

func (st *MemStorage) Load() error {
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

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		logger.Log.Error("close file error", zap.Error(err))
	}
}
