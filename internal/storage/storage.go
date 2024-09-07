package storage

import (
	"strconv"
)

type Gauge = float64
type Counter = int64

type MemStorage struct {
	GaugeMap   map[string]Gauge
	CounterMap map[string][]Counter
}

func (storage *MemStorage) AddCounter(name string, item string) error {
	counter, err := strconv.ParseInt(item, 10, 64)
	if err == nil {
		storage.CounterMap[name] = append(storage.CounterMap[name], counter)
	}
	return err

}

func (storage *MemStorage) AddGauge(name string, item string) error {
	gauge, err := strconv.ParseFloat(item, 64)
	if err == nil {
		storage.GaugeMap[name] = gauge
	}
	return err

}
