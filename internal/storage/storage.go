package storage

type Gauge = float64
type Counter = int64

type MemStorage struct {
	GaugeMap   map[string]Gauge
	CounterMap map[string][]Counter
}

func (storage *MemStorage) AddCounter(name string, item Counter) {
	storage.CounterMap[name] = append(storage.CounterMap[name], item)
}

func (storage *MemStorage) AddGauge(name string, item Gauge) {
	storage.GaugeMap[name] = item
}
