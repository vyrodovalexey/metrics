package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"log"
)

type Config struct {
	ListenAddr string `env:"ADDRESS"`
}

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}
	if len(cfg.ListenAddr) == 0 {
		flag.StringVar(&cfg.ListenAddr, "a", ":8080", "input ip:port to listen")
		flag.Parse()
	}

	r := SetupRouter(&mst)
	r.LoadHTMLGlob("templates/*")
	r.Run(cfg.ListenAddr)

}
