package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"go.uber.org/zap"
)

type PGStorage struct {
	db *sql.DB
}

func NewPGStorage(db *sql.DB) *PGStorage {
	return &PGStorage{
		db: db,
	}
}

func (st *PGStorage) UpdateGauge(ctx context.Context, metric models.Metric) error {
	_, err := st.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2)
		ON CONFLICT (id, type) DO UPDATE SET value = $2`, metric.ID, metric.Value)
	if err != nil {
		return err
	}
	return nil
}

func (st *PGStorage) GetGauge(ctx context.Context, id string) (models.Metric, error) {
	row := st.db.QueryRowContext(ctx, `SELECT id, type, value FROM metrics WHERE id = $1 AND type = 'gauge'`, id)
	var metric models.Metric
	err := row.Scan(&metric.ID, &metric.MType, &metric.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metric{}, storage.ErrNotFound
		}
		return models.Metric{}, err
	}

	return metric, nil
}

func (st *PGStorage) GetGaugeList(ctx context.Context) storage.GaugeList {
	rows, err := st.db.QueryContext(ctx, `SELECT id, type, value FROM metrics WHERE type = 'gauge'`)
	if err != nil {
		logger.Log.Error("failed to get gauge list", zap.Error(err))
		return nil
	}
	defer rows.Close()

	gauges := make(storage.GaugeList)
	for rows.Next() {
		var metric models.Metric
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Value)
		if err != nil {
			logger.Log.Error("failed to scan gauge", zap.Error(err))
			continue
		}
		gauges[metric.ID] = metric
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("failed to iterate over gauges", zap.Error(err))
	}

	return gauges
}

func (st *PGStorage) UpdateCounter(ctx context.Context, metric models.Metric) error {
	_, err := st.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'counter', $2)
		ON CONFLICT (id, type) DO UPDATE SET value = metrics.value + excluded.value`, metric.ID, metric.Delta)
	if err != nil {
		return err
	}
	return nil
}

func (st *PGStorage) GetCounter(ctx context.Context, id string) (models.Metric, error) {
	row := st.db.QueryRowContext(ctx, `SELECT id, type, value FROM metrics WHERE id = $1 AND type = 'counter'`, id)
	var metric models.Metric
	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metric{}, storage.ErrNotFound
		}
		return models.Metric{}, err
	}

	return metric, nil
}

func (st *PGStorage) GetCounterList(ctx context.Context) storage.CounterList {
	rows, err := st.db.QueryContext(ctx, `SELECT id, type, value FROM metrics WHERE type = 'counter'`)
	if err != nil {
		logger.Log.Error("failed to get gauge list", zap.Error(err))
		return nil
	}
	defer rows.Close()

	counters := make(storage.CounterList)
	for rows.Next() {
		var metric models.Metric
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta)
		if err != nil {
			logger.Log.Error("failed to scan counter", zap.Error(err))
			continue
		}
		counters[metric.ID] = metric
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("failed to iterate over counters", zap.Error(err))
	}

	return counters
}

func (st *PGStorage) Save(ctx context.Context) error {
	return nil
}

func (st *PGStorage) Load(ctx context.Context) error {
	return nil
}

func (st *PGStorage) Ping(ctx context.Context) error {
	return st.db.Ping()
}

func (st *PGStorage) Close(ctx context.Context) {
	if err := st.db.Close(); err != nil {
		logger.Log.Fatal("failed to close database connection", zap.Error(err))
	}
}
