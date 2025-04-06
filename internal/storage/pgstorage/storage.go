package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"time"
)

type PGStorage struct {
	db *sql.DB
}

func NewPGStorage(db *sql.DB) *PGStorage {
	return &PGStorage{
		db: db,
	}
}

func isRetriableError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ConnectionException, pgerrcode.ConnectionDoesNotExist, pgerrcode.ConnectionFailure:
			return true
		}
	}
	return false
}

func withRetries(fn func() error) error {
	var err error
	for i, delay := range []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second} {
		err = fn()
		if err == nil || !isRetriableError(err) {
			break
		}
		logger.Log.Warn("retriable error occurred, retrying...", zap.Error(err), zap.Int("attempt", i+1))
		time.Sleep(delay)
	}
	return err
}

func (st *PGStorage) UpdateGauge(ctx context.Context, metric models.Metric) error {
	return withRetries(func() error {
		_, err := st.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2)
		ON CONFLICT (id, type) DO UPDATE SET value = $2`, metric.ID, metric.Value)
		return err
	})
}

func (st *PGStorage) GetGauge(ctx context.Context, id string) (models.Metric, error) {
	var metric models.Metric
	err := withRetries(func() error {
		row := st.db.QueryRowContext(ctx, `SELECT id, type, value FROM metrics WHERE id = $1 AND type = 'gauge'`, id)
		return row.Scan(&metric.ID, &metric.MType, &metric.Value)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metric{}, storage.ErrNotFound
		}
		return models.Metric{}, err
	}
	return metric, nil
}

func (st *PGStorage) GetGaugeList(ctx context.Context) storage.GaugeList {
	var gauges storage.GaugeList
	err := withRetries(func() error {
		rows, err := st.db.QueryContext(ctx, `SELECT id, type, value FROM metrics WHERE type = 'gauge'`)
		if err != nil {
			return err
		}
		defer rows.Close()

		gauges = make(storage.GaugeList)
		for rows.Next() {
			var metric models.Metric
			err := rows.Scan(&metric.ID, &metric.MType, &metric.Value)
			if err != nil {
				logger.Log.Error("failed to scan gauge", zap.Error(err))
				continue
			}
			gauges[metric.ID] = metric
		}
		return rows.Err()
	})
	if err != nil {
		logger.Log.Error("failed to get gauge list", zap.Error(err))
	}
	return gauges
}

func (st *PGStorage) UpdateCounter(ctx context.Context, metric models.Metric) error {
	return withRetries(func() error {
		_, err := st.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'counter', $2)
			ON CONFLICT (id, type) DO UPDATE SET value = metrics.value + excluded.value`, metric.ID, metric.Delta)
		return err
	})
}

func (st *PGStorage) GetCounter(ctx context.Context, id string) (models.Metric, error) {
	var metric models.Metric
	var value float64
	err := withRetries(func() error {
		row := st.db.QueryRowContext(ctx, `SELECT id, type, value FROM metrics WHERE id = $1 AND type = 'counter'`, id)
		return row.Scan(&metric.ID, &metric.MType, &value)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metric{}, storage.ErrNotFound
		}
		return models.Metric{}, err
	}
	delta := int64(value)
	metric.Delta = &delta
	return metric, nil
}

func (st *PGStorage) GetCounterList(ctx context.Context) storage.CounterList {
	var counters storage.CounterList
	err := withRetries(func() error {
		rows, err := st.db.QueryContext(ctx, `SELECT id, type, value FROM metrics WHERE type = 'counter'`)
		if err != nil {
			return err
		}
		defer rows.Close()

		counters = make(storage.CounterList)
		for rows.Next() {
			var metric models.Metric
			err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta)
			if err != nil {
				logger.Log.Error("failed to scan counter", zap.Error(err))
				continue
			}
			counters[metric.ID] = metric
		}
		return rows.Err()
	})
	if err != nil {
		logger.Log.Error("failed to get counter list", zap.Error(err))
	}
	return counters
}

func (st *PGStorage) UpdateBatch(ctx context.Context, metrics []models.Metric) error {
	return withRetries(func() error {
		tx, err := st.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer func(tx *sql.Tx) {
			err = tx.Rollback()
			if err != nil && !errors.Is(err, sql.ErrTxDone) {
				logger.Log.Error("failed to rollback transaction", zap.Error(err))
			}
		}(tx)

		for _, metric := range metrics {
			switch metric.MType {
			case models.TypeGauge:
				_, err = tx.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2)
     ON CONFLICT (id, type) DO UPDATE SET value = $2`, metric.ID, metric.Value)
				if err != nil {
					return err
				}
			case models.TypeCounter:
				_, err = tx.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'counter', $2)
     ON CONFLICT (id, type) DO UPDATE SET value = metrics.value + excluded.value`, metric.ID, metric.Delta)
				if err != nil {
					return err
				}
			default:
				return storage.ErrWrongType
			}
		}

		return tx.Commit()
	})
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
