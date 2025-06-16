package storage

import (
	"context"
	"errors"

	"github.com/invinciblewest/metrics/internal/models"
)

type GaugeList map[string]models.Metric   // GaugeList содержит метрики типа Gauge, где ключ - это идентификатор метрики, а значение - сама метрика.
type CounterList map[string]models.Metric // CounterList содержит метрики типа Counter, где ключ - это идентификатор метрики, а значение - сама метрика.

var (
	ErrNotFound  = errors.New("not found")
	ErrWrongType = errors.New("wrong type")
)

// Storage интерфейс для работы с хранилищем метрик.
type Storage interface {
	UpdateGauge(ctx context.Context, metric models.Metric) error      // UpdateGauge обновляет метрику типа Gauge в хранилище.
	GetGauge(ctx context.Context, id string) (models.Metric, error)   // GetGauge извлекает метрику типа Gauge из хранилища по идентификатору.
	GetGaugeList(ctx context.Context) GaugeList                       // GetGaugeList возвращает список всех метрик типа Gauge в хранилище.
	UpdateCounter(ctx context.Context, metric models.Metric) error    // UpdateCounter обновляет метрику типа Counter в хранилище.
	GetCounter(ctx context.Context, id string) (models.Metric, error) // GetCounter извлекает метрику типа Counter из хранилища по идентификатору.
	GetCounterList(ctx context.Context) CounterList                   // GetCounterList возвращает список всех метрик типа Counter в хранилище.
	UpdateBatch(ctx context.Context, metrics []models.Metric) error   // UpdateBatch обновляет пакет метрик в хранилище.
	Save(ctx context.Context) error                                   // Save сохраняет текущее состояние хранилища в постоянное хранилище (например, файл или базу данных).
	Load(ctx context.Context) error                                   // Load загружает состояние хранилища из постоянного хранилища (например, файла или базы данных).
	Ping(ctx context.Context) error                                   // Ping проверяет доступность хранилища.
	Close(ctx context.Context) error                                  // Close закрывает соединение с хранилищем и освобождает ресурсы.
}
