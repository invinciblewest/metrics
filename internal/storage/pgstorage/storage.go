package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/avast/retry-go"
	"github.com/invinciblewest/metrics/internal/logger"
	"github.com/invinciblewest/metrics/internal/models"
	"github.com/invinciblewest/metrics/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

// PGStorage представляет собой хранилище метрик в PostgreSQL.
type PGStorage struct {
	db *sql.DB
}

// NewPGStorage создает новый экземпляр PGStorage с заданным подключением к базе данных.
func NewPGStorage(db *sql.DB) *PGStorage {
	return &PGStorage{
		db: db,
	}
}

// isRetriableError проверяет, является ли ошибка временной и может быть повторена.
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

// withRetries выполняет функцию с повторными попытками в случае временных ошибок.
func withRetries(ctx context.Context, fn func() error) error {
	retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	return retry.Do(
		fn,
		retry.Context(ctx),
		retry.Attempts(3),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if n < uint(len(retryDelays)) {
				return retryDelays[n-1]
			}
			return 0
		}),
		retry.RetryIf(isRetriableError),
		retry.OnRetry(func(n uint, err error) {
			if err != nil {
				logger.Log.Warn("retriable error occurred, retrying...", zap.Error(err), zap.Uint("attempt", n+1))
			}
		}),
	)
}

// UpdateGauge обновляет метрику типа Gauge в хранилище.
func (st *PGStorage) UpdateGauge(ctx context.Context, metric models.Metric) error {
	return withRetries(ctx, func() error {
		_, err := st.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2)
		ON CONFLICT (id, type) DO UPDATE SET value = $2`, metric.ID, metric.Value)
		return err
	})
}

// GetGauge извлекает метрику типа Gauge из хранилища по идентификатору.
func (st *PGStorage) GetGauge(ctx context.Context, id string) (models.Metric, error) {
	var metric models.Metric
	err := withRetries(ctx, func() error {
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

// GetGaugeList возвращает список всех метрик типа Gauge в хранилище.
func (st *PGStorage) GetGaugeList(ctx context.Context) storage.GaugeList {
	var gauges storage.GaugeList
	err := withRetries(ctx, func() error {
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

// UpdateCounter обновляет метрику типа Counter в хранилище.
func (st *PGStorage) UpdateCounter(ctx context.Context, metric models.Metric) error {
	return withRetries(ctx, func() error {
		_, err := st.db.ExecContext(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, 'counter', $2)
			ON CONFLICT (id, type) DO UPDATE SET value = metrics.value + excluded.value`, metric.ID, metric.Delta)
		return err
	})
}

// GetCounter извлекает метрику типа Counter из хранилища по идентификатору.
func (st *PGStorage) GetCounter(ctx context.Context, id string) (models.Metric, error) {
	var metric models.Metric
	var value float64
	err := withRetries(ctx, func() error {
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

// GetCounterList возвращает список всех метрик типа Counter в хранилище.
func (st *PGStorage) GetCounterList(ctx context.Context) storage.CounterList {
	var counters storage.CounterList
	err := withRetries(ctx, func() error {
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

// UpdateBatch обновляет пакет метрик в хранилище.
func (st *PGStorage) UpdateBatch(ctx context.Context, metrics []models.Metric) error {
	return withRetries(ctx, func() error {
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

// Save сохраняет текущее состояние хранилища в постоянное хранилище.
// В данном случае, сохранение в PostgreSQL не требуется, так как все изменения уже сохраняются в базе данных.
func (st *PGStorage) Save(ctx context.Context) error {
	return nil
}

// Load загружает состояние хранилища из постоянного хранилища.
// В данном случае, загрузка из PostgreSQL не требуется, так как все данные уже находятся в базе данных.
func (st *PGStorage) Load(ctx context.Context) error {
	return nil
}

// Ping проверяет доступность хранилища, отправляя простой запрос к базе данных.
func (st *PGStorage) Ping(ctx context.Context) error {
	return st.db.Ping()
}

// Close закрывает соединение с хранилищем и освобождает ресурсы.
func (st *PGStorage) Close(ctx context.Context) error {
	return st.db.Close()
}
