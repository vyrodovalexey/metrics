package main

import (
	"github.com/vyrodovalexey/metrics/internal/handlers"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"log"
	"net/http"
)

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}

	http.HandleFunc("/update/", handlers.Update(&mst))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
