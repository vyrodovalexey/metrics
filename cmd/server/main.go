package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

const (
	serverAddr = ":8080"
)

type Config struct {
	ListenAddr string `env:"ADDRESS"`
}

func main() {

	gauge := make(map[string]storage.Gauge)
	counter := make(map[string]storage.Counter)
	mst := storage.MemStorage{GaugeMap: gauge, CounterMap: counter}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.DisableCaller = true

	logger, logerr := loggerConfig.Build()
	if logerr != nil {
		log.Fatalf("can't initialize zap logger: %v", logerr)
	}
	// nolint:errcheck
	defer logger.Sync()

	sugar := logger.Sugar()

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}
	if len(cfg.ListenAddr) == 0 {
		flag.StringVar(&cfg.ListenAddr, "a", serverAddr, "input ip:port to listen")
		flag.Parse()
	}

	r := SetupRouter(&mst, sugar)
	r.LoadHTMLGlob("templates/*")
	r.Run(cfg.ListenAddr)

}
