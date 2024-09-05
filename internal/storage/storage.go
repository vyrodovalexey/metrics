package storage

import "time"

type gauge = float64
type counter = int64

type GaugeItem struct {
	Value gauge
	Date  time.Time
}

type CounterItem struct {
	Value counter
	Date  time.Time
}

type MemStorage struct {
	Gauge   map[string]GaugeItem
	Counter map[string][]CounterItem
}

func (storage *MemStorage) AddCounter(name string, item CounterItem) {
	storage.Counter[name] = append(storage.Counter[name], item)
}

func (storage *MemStorage) AddGauge(name string, item GaugeItem) {
	storage.Gauge[name] = item
}

var Storage MemStorage
