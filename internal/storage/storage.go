package storage

type Gauge = float64
type Counter = int64

type Storage interface {
	Init()
	AddGaugeAsString(key string, value string) error
	AddGauge(key string, value Gauge)
	GetGauge(key string) (Gauge, bool)
	AddCounterAsString(key string, value string) error
	AddCounter(key string, value Counter)
	GetCounter(key string) (Counter, bool)
	GetAllMetricNames() (map[string]string, map[string]string)
}
