package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	pathSlice := strings.Split(r.URL.Path[1:], "/")

	if len(pathSlice) == 3 && (pathSlice[1] == "gauge" || pathSlice[1] == "counter") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(pathSlice) != 4 || (pathSlice[1] != "gauge" && pathSlice[1] != "counter") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pathSlice[1] == "gauge" {
		_, err := strconv.ParseFloat(pathSlice[3], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			gauge, _ := strconv.ParseFloat(pathSlice[3], 64)
			Storage.AddGauge(pathSlice[2], GaugeItem{gauge, time.Now()})
		}
	}

	if pathSlice[1] == "counter" {
		_, err := strconv.ParseInt(pathSlice[3], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			counter, _ := strconv.ParseInt(pathSlice[3], 10, 64)

			Storage.AddCounter(pathSlice[2], CounterItem{counter, time.Now()})
		}
	}

}

func run() error {
	return http.ListenAndServe(`:8080`, nil)
}

func main() {
	gauge := make(map[string]GaugeItem)
	counter := make(map[string][]CounterItem)

	Storage = MemStorage{gauge, counter}

	http.HandleFunc("/update/", Update)
	if err := run(); err != nil {
		panic(err)
	}
}
