package main

import (
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
)

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string][]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}

	mux := http.NewServeMux()
	wrappedMuxUpdate := handlers.NewStorageHandler(mux, &mst)
	mux.HandleFunc("/update/", handlers.Update)
	err := http.ListenAndServe(`:8080`, wrappedMuxUpdate)

	if err != nil {
		panic(err)
	}
}
