package storage

import (
	"os"
)

type Gauge = float64
type Counter = int64

type Storage interface {
	// New Создание нового хранилища
	New()
	// Load Загрузка хранилища из файла
	Load(f *os.File) error
	// AddGaugeAsString Добавление метрики Gauge с строковым значением
	AddGaugeAsString(key string, value string) error
	// AddGauge Добавление метрики Gauge с числовым значением
	AddGauge(key string, value Gauge)
	// GetGauge Получение метрики Gauge по ключу
	GetGauge(key string) (Gauge, bool)
	// AddCounterAsString Добавление метрики Counter с строковым значением
	AddCounterAsString(key string, value string) error
	// AddCounter Добавление метрики Counter
	AddCounter(key string, value Counter)
	// GetCounter Получение метрики Counter  по ключу
	GetCounter(key string) (Counter, bool)
	// GetAllMetricNames Получение списка имен метрик
	GetAllMetricNames() (map[string]string, map[string]string)
	// SaveAsync Асинхронная сохранение данных хранилища в файл
	SaveAsync(f *os.File, interval int)
	// Save Сохранение данных хранилища в файл
	Save(f *os.File)
}
