package storage

import (
	"github.com/vyrodovalexey/metrics/internal/model"
)

type Storage interface {
	// New Создание нового хранилища
	NewMemStorage(filePath string, interval uint) error
	// NewDatabaseConnection Создание соединения с базой данных
	//NewDatabaseConnection(c string) (pgx.Conn, error)
	// CheckDatabaseConnection Проверка соединения с базой данных
	//CheckDatabaseConnection() error
	// Load Загрузка хранилища из файла
	LoadMemStorage(filePath string, interval uint) error
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
	GetAllMetricNames() (map[string]string, map[string]string)
	// SaveAsync Асинхронная сохранение данных хранилища в файл
	SaveAsync() error
	// Save Сохранение данных хранилища в файл
	Save() error

	Close()
}
