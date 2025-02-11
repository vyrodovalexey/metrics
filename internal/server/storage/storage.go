package storage

import (
	"github.com/vyrodovalexey/metrics/internal/model"
	"go.uber.org/zap"
)

type Storage interface {
	// New Создание нового хранилища
	New(filePath string, interval uint, log *zap.SugaredLogger) error
	// NewDatabaseConnection Создание соединения с базой данных
	Check() error
	// Load Загрузка хранилища из файла
	Load(filePath string, interval uint, log *zap.SugaredLogger) error
	// UpdateGauge Добавление метрики Gauge
	UpdateGauge(name string, item model.Gauge) error
	// UpdateCounter Добавление метрики Counter
	UpdateCounter(name string, item model.Counter) error
	// UpdateMetric Добавление метрики
	UpdateMetric(metrics *model.Metrics) error
	// GetGauge Получение метрики Gauge
	GetGauge(name string) (model.Gauge, bool)
	// GetCounter Получение метрики Counter
	GetCounter(name string) (model.Counter, bool)
	// GetMetric Получение метрики
	GetMetric(metrics *model.Metrics) bool
	// GetAllMetricNames Получение списка имен метрик
	GetAllMetricNames() (map[string]string, map[string]string, error)
	// SaveAsync Асинхронная сохранение данных хранилища в файл
	SaveAsync()
	// Save Сохранение данных хранилища в файл
	Save() error

	Close()
}
