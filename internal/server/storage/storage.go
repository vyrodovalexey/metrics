package storage

import (
	"context"
	"github.com/vyrodovalexey/metrics/internal/model"
)

type Storage interface {
	// New Создание нового хранилища
	New(ctx context.Context, filePath string, interval uint) error
	// NewDatabaseConnection Создание соединения с базой данных
	Check(ctx context.Context) error
	// Load Загрузка хранилища из файла
	Load(ctx context.Context, filePath string, interval uint) error
	// UpdateGauge Добавление метрики Gauge
	UpdateGauge(ctx context.Context, name string, item model.Gauge) error
	// UpdateCounter Добавление метрики Counter
	UpdateCounter(ctx context.Context, name string, item model.Counter) error
	// UpdateMetric Добавление метрики
	UpdateMetric(ctx context.Context, metrics *model.Metrics) error
	// GetGauge Получение метрики Gauge
	GetGauge(ctx context.Context, name string) (model.Gauge, bool)
	// GetCounter Получение метрики Counter
	GetCounter(ctx context.Context, name string) (model.Counter, bool)
	// GetMetric Получение метрики
	GetMetric(ctx context.Context, metrics *model.Metrics) bool
	// GetAllMetricNames Получение списка имен метрик
	GetAllMetricNames(ctx context.Context) (map[string]string, map[string]string, error)
	// SaveAsync Асинхронная сохранение данных хранилища в файл
	SaveAsync() error
	// Save Сохранение данных хранилища в файл
	Save() error

	Close()
}
