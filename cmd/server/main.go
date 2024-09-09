package main

import (
	"github.com/vyrodovalexey/metrics/internal/storage"
)

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string][]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}

	r := SetupRouter(&mst)
	r.LoadHTMLGlob("templates/*")
	r.Run(":8080")

}
