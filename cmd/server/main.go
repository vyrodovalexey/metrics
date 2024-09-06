package main

import (
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Storage storage.MemStorage

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
			Storage.AddGauge(pathSlice[2], storage.GaugeItem{gauge, time.Now()})
		}
	}

	if pathSlice[1] == "counter" {
		_, err := strconv.ParseInt(pathSlice[3], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			counter, _ := strconv.ParseInt(pathSlice[3], 10, 64)

			Storage.AddCounter(pathSlice[2], storage.CounterItem{counter, time.Now()})
		}
	}

}

func run() error {
	return http.ListenAndServe(`:8080`, nil)
}

func main() {
	gauge := make(map[string]storage.GaugeItem)
	counter := make(map[string][]storage.CounterItem)

	Storage = storage.MemStorage{gauge, counter}

	http.HandleFunc("/update/", Update)
	if err := run(); err != nil {
		panic(err)
	}
}
