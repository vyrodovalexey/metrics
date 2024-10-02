package storage

import (
	"github.com/vyrodovalexey/metrics/internal/model"
	"os"
)

type Gauge = float64
type Counter = int64

type Storage interface {
	// New Создание нового хранилища
	New()
	// Load Загрузка хранилища из файла
	Load(f *os.File) error
	// UpdateGauge Добавление метрики Gauge
	UpdateGauge(name string, item Gauge, f *os.File, p bool) error
	// UpdateCounter Добавление метрики Counter
	UpdateCounter(name string, item Counter, f *os.File, p bool) error
	// UpdateMetric Добавление метрики
	UpdateMetric(metrics *model.Metrics, f *os.File, p bool) error
	// GetGauge Получение метрики Gauge
	GetGauge(name string) (Gauge, bool)
	// GetCounter Получение метрики Counter
	GetCounter(name string) (Counter, bool)
	// GetMetric Получение метрики
	GetMetric(metrics *model.Metrics) bool
	// GetAllMetricNames Получение списка имен метрик
	GetAllMetricNames() (map[string]string, map[string]string)
	// SaveAsync Асинхронная сохранение данных хранилища в файл
	SaveAsync(f *os.File, interval uint)
	// Save Сохранение данных хранилища в файл
	Save(f *os.File) error
}
