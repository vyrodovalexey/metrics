package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

const (
	defaultListenAddr    = ":8080"
	defaultStoreInterval = 300
	defaultFileStorePath = "storage.json"
	defaultRestore       = true
)

type Config struct {
	ListenAddr    string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FileStorePath string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
}

func main() {

	var st storage.Storage = &storage.MemStorage{}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.DisableCaller = true

	logger, logerr := loggerConfig.Build()
	if logerr != nil {
		log.Fatalf("can't initialize zap logger: %v", logerr)
	}
	// nolint:errcheck
	defer logger.Sync()

	sugarLog := logger.Sugar()

	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {

		sugarLog.DPanicf("can't parse config: %v", err)
	}
	if len(cfg.ListenAddr) == 0 {
		flag.StringVar(&cfg.ListenAddr, "a", defaultListenAddr, "input ip:port to listen")
	}
	if cfg.StoreInterval < 1 {
		flag.IntVar(&cfg.StoreInterval, "i", defaultStoreInterval, "seconds delay between save data to file ")
	}
	if len(cfg.FileStorePath) == 0 {
		flag.StringVar(&cfg.FileStorePath, "f", defaultFileStorePath, "path to file")
	}
	if !cfg.Restore {
		flag.BoolVar(&cfg.Restore, "r", defaultRestore, "restore date on load")
	}
	flag.Parse()

	sugarLog.Infow("Server starting with",
		"address", cfg.ListenAddr,
		"File store path", cfg.FileStorePath,
		"Load file true/false", cfg.Restore,
		"Store interval in sec", cfg.StoreInterval,
	)

	file, ferr := os.OpenFile(cfg.FileStorePath, os.O_RDWR|os.O_CREATE, 0666)
	if ferr != nil {
		sugarLog.Infow("Error creating file:",
			"error", ferr,
		)
		return
	}

	if cfg.Restore {
		st.Load(file)
	} else {
		st.New()
	}

	if cfg.StoreInterval > 0 {
		go st.SaveAsync(file, cfg.StoreInterval)
	}
	r := SetupRouter(st, sugarLog)
	r.LoadHTMLGlob("templates/*")
	r.Run(cfg.ListenAddr)
	st.Save(file)
	defer file.Close()
}
