package pgstorage

import (
	"database/sql"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
)

type PGStorage struct {
	db *sql.DB
}

func NewPGStorage(dsn string) *PGStorage {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	return &PGStorage{
		db: db,
	}
}

func (st *PGStorage) UpdateGauge(metric models.Metric) error {
	return nil
}

func (st *PGStorage) GetGauge(id string) (models.Metric, error) {
	return models.Metric{}, nil
}

func (st *PGStorage) GetGaugeList() storage.GaugeList {
	return nil
}

func (st *PGStorage) UpdateCounter(metric models.Metric) error {
	return nil
}

func (st *PGStorage) GetCounter(id string) (models.Metric, error) {
	return models.Metric{}, nil
}

func (st *PGStorage) GetCounterList() storage.CounterList {
	return nil
}

func (st *PGStorage) Save() error {
	return nil
}

func (st *PGStorage) Load() error {
	return nil
}

func (st *PGStorage) Ping() error {
	return st.db.Ping()
}

func (st *PGStorage) Close() {
	if err := st.db.Close(); err != nil {
		logger.Log.Fatal("failed to close database connection", zap.Error(err))
	}
}
