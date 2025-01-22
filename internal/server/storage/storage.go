package storage

import (
	"github.com/vyrodovalexey/metrics/internal/model"
	"os"
)

type Storage interface {
	// New Создание нового хранилища
	New()
	// NewDatabaseConnection Создание соединения с базой данных
	//NewDatabaseConnection(c string) (pgx.Conn, error)
	// CheckDatabaseConnection Проверка соединения с базой данных
	//CheckDatabaseConnection() error
	// Load Загрузка хранилища из файла
	Load(f *os.File) error
	// UpdateGauge Добавление метрики Gauge
	UpdateGauge(name string, item model.Gauge, f *os.File, p bool) error
	// UpdateCounter Добавление метрики Counter
	UpdateCounter(name string, item model.Counter, f *os.File, p bool) error
	// UpdateMetric Добавление метрики
	UpdateMetric(metrics *model.Metrics, f *os.File, p bool) error
	// GetGauge Получение метрики Gauge
	GetGauge(name string) (model.Gauge, bool)
	// GetCounter Получение метрики Counter
	GetCounter(name string) (model.Counter, bool)
	// GetMetric Получение метрики
	GetMetric(metrics *model.Metrics) bool
	// GetAllMetricNames Получение списка имен метрик
	GetAllMetricNames() (map[string]string, map[string]string)
	// SaveAsync Асинхронная сохранение данных хранилища в файл
	SaveAsync(f *os.File, interval uint) error
	// Save Сохранение данных хранилища в файл
	Save(f *os.File) error
}
