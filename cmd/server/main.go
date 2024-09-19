package main

import (
	"flag"
	"github.com/vyrodovalexey/metrics/internal/storage"
)

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}
	listenAddr := flag.String("a", ":8080", "input ip:port to listen")

	flag.Parse()

	r := SetupRouter(&mst)
	r.LoadHTMLGlob("templates/*")
	r.Run(*listenAddr)

}
