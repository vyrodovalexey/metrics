package storage

import (
	"fmt"
	"strconv"
)

type Gauge = float64
type Counter = int64

type MemStorage struct {
	GaugeMap   map[string]Gauge
	CounterMap map[string]Counter
}

func (storage *MemStorage) AddCounter(name string, item string) error {
	counter, err := strconv.ParseInt(item, 10, 64)
	if err == nil {
		storage.CounterMap[name] = storage.CounterMap[name] + counter
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

func (storage *MemStorage) GetAllMetricNames() (map[string]string, map[string]string) {
	//to debug
	//names := make([]string, 0, len(storage.GaugeMap)+len(storage.CounterMap))
	gvalues := make(map[string]string, len(storage.GaugeMap))

	cvalues := make(map[string]string, len(storage.CounterMap))
	// Iterate over the map and collect the keys
	for name := range storage.GaugeMap {
		//to debug
		//names = append(names, name)
		gv, _ := storage.GetGauge(name)
		gvalues[name] = fmt.Sprintf("%v", gv)
	}

	for name := range storage.CounterMap {
		//to debug
		//names = append(names, name)
		cv, _ := storage.GetCounter(name)
		cvalues[name] = fmt.Sprintf("%v", cv)
	}

	return gvalues, cvalues
}

func (storage *MemStorage) GetGauge(name string) (Gauge, bool) {
	res, e := storage.GaugeMap[name]
	if e {
		return res, e
	}
	return 0, false

}

func (storage *MemStorage) GetCounter(name string) (Counter, bool) {
	res, e := storage.CounterMap[name]

	if e {
		return res, e
	}
	return 0, false
}
