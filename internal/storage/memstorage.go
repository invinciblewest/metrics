package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"go.uber.org/zap"
	"os"
)

type MemStorage struct {
	Gauges   GaugeList   `json:"gauges"`
	Counters CounterList `json:"counters"`
	path     string
	syncSave bool
}

func NewMemStorage(path string, syncSave bool) *MemStorage {
	return &MemStorage{
		Gauges:   make(GaugeList),
		Counters: make(CounterList),
		path:     path,
		syncSave: syncSave,
	}
}

func (st *MemStorage) UpdateGauge(metrics models.Metrics) error {
	if metrics.MType != models.TypeGauge {
		return errors.New("wrong type")
	}

	st.Gauges[metrics.ID] = metrics
	if st.syncSave {
		return st.Save()
	}
	return nil
}

func (st *MemStorage) GetGauge(id string) (models.Metrics, error) {
	value, exists := st.Gauges[id]
	if exists {
		return value, nil
	} else {
		return models.Metrics{}, errors.New("not found")
	}
}

func (st *MemStorage) GetGaugeList() GaugeList {
	return st.Gauges
}

func (st *MemStorage) UpdateCounter(metrics models.Metrics) error {
	if metrics.MType != models.TypeCounter {
		return errors.New("wrong type")
	}

	currentMetrics, exists := st.Counters[metrics.ID]
	if exists {
		*metrics.Delta += *currentMetrics.Delta
	}

	st.Counters[metrics.ID] = metrics

	if st.syncSave {
		return st.Save()
	}
	return nil
}

func (st *MemStorage) GetCounter(id string) (models.Metrics, error) {
	value, exists := st.Counters[id]
	if exists {
		return value, nil
	} else {
		return models.Metrics{}, errors.New("not found")
	}
}

func (st *MemStorage) GetCounterList() CounterList {
	return st.Counters
}

func (st *MemStorage) Save() error {
	if st.path == "" {
		return nil
	}
	logger.Log.Info("saving storage...", zap.String("storage", st.path))

	file, err := os.OpenFile(st.path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer closeFile(file)

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(*st); err != nil { // Передаем саму структуру
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

	file, err := os.OpenFile(st.path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer closeFile(file)

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
