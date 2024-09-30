package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"log"
)

func ConfigParser(cfg *config.Config) {

	flag.StringVar(&cfg.EndpointAddr, "a", cfg.EndpointAddr, "input ip:port or host:port of metrics server")

	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "seconds delay interval to send metrics to metrics server")
	flag.IntVar(&cfg.PoolInterval, "p", cfg.PoolInterval, "seconds delay between scribing metrics from host")

	flag.Parse()

	err := env.Parse(cfg)

	if err != nil {
		log.Fatalf("can't parse ENV: %v", err)
	}

}
