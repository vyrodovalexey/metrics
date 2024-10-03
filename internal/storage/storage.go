package storage

type Gauge = float64
type Counter = int64

type Storage interface {
	Init()
	AddGauge(key string, value string) error
	GetGauge(key string) (Gauge, bool)
	AddCounter(key string, value string) error
	GetCounter(key string) (Counter, bool)
	GetAllMetricNames() (map[string]string, map[string]string)
}
