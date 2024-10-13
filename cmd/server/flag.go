package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/vyrodovalexey/metrics/internal/server/config"
	"log"
)

func ConfigParser(cfg *config.Config) {

	// Устанавливаем флаг для адреса прослушивания (IP:порт)
	flag.StringVar(&cfg.ListenAddr, "a", cfg.ListenAddr, "input ip:port to listen")
	// Устанавливаем флаг для интервала сохранения данных в секундах
	flag.UintVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "seconds delay between save data to file ")
	// Устанавливаем флаг для пути к файлу хранения
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "path to file")
	// Устанавливаем флаг для восстановления данных при загрузке
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore date on load")
	// Устанавливаем флаг для адреса базы данных
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database connection string")
	flag.Parse() // Парсим флаги командной строки

	// Парсим переменные окружения и сохраняем их в конфигурацию и перезаписывая существующие
	err := env.Parse(cfg)

	if err != nil {
		// Логируем ошибку, если не удалось распарсить переменные окружения
		log.Fatalf("can't parse ENV: %v", err)
	}

}
