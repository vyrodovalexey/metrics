package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/agent/config"
	"log"
)

func ConfigParser(cfg *config.Config) {
	// Устанавливаем флаг для адреса сервера метрик (IP:порт или хост:порт)
	flag.StringVar(&cfg.EndpointAddr, "a", cfg.EndpointAddr, "input ip:port or host:port of metrics server")
	// Устанавливаем флаг для интервала отправки метрик в секундах
	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "seconds delay interval to send metrics to metrics server")
	// Устанавливаем флаг для интервала опроса метрик в секундах
	flag.IntVar(&cfg.PoolInterval, "p", cfg.PoolInterval, "seconds delay between scribing metrics from host")

	// Парсим флаги командной строки
	flag.Parse()

	// Парсим переменные окружения и сохраняем их в конфигурацию и перезаписывая существующие
	err := env.Parse(cfg)

	if err != nil {
		// Логируем ошибку, если не удалось распарсить переменные окружения
		log.Fatalf("can't parse ENV: %v", err)
	}

}
