package storage

import (
	"os"
)

type Gauge = float64
type Counter = int64

type Storage interface {
	New()
	Load(f *os.File) error
	AddGaugeAsString(key string, value string) error
	AddGauge(key string, value Gauge)
	GetGauge(key string) (Gauge, bool)
	AddCounterAsString(key string, value string) error
	AddCounter(key string, value Counter)
	GetCounter(key string) (Counter, bool)
	GetAllMetricNames() (map[string]string, map[string]string)
	SaveAsync(f *os.File, interval int)
	Save(f *os.File)
}
