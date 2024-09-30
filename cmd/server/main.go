package main

import (
	"github.com/vyrodovalexey/metrics/internal/server/config"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"go.uber.org/zap"
	"os"
)

func main() {
	// Создаем новый экземпляр конфигурации
	cfg := config.New()
	// Парсим настройки конфигурации
	ConfigParser(cfg)
	// Инициализируем логирование с уровнем Info
	lg := logging.NewLogging(zap.InfoLevel)

	// Логируем информацию о запуске сервера
	lg.Infow("Server starting with",
		"address", cfg.ListenAddr,
		"File store path", cfg.FileStoragePath,
		"Load storage file on start true/false", cfg.Restore,
		"Store interval in sec", cfg.StoreInterval,
	)
	// Инициализируем интерфейс и структуру хранения данных
	var st storage.Storage = &storage.MemStorage{}

	// Открываем или создаем файл для хранения
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		// Логируем ошибку, если открытие/создание файла не удалось
		lg.Infow("Initializing file storage...",
			"Error creating file:", err,
		)
		return
	}
	lg.Infow("File storage initialized")

	// Проверяем, нужно ли загружать файл харнилища
	// Если нет, инициализируем новое
	if cfg.Restore {
		st.Load(file)
	} else {
		st.New()
	}

	// При ассинхронном режиме запускаем фоновое сохранение структуры данных
	if cfg.StoreInterval > 0 {
		go st.SaveAsync(file, cfg.StoreInterval)
	}
	// Инициализируем маршрутизатор с хранилищем и логированием
	r := SetupRouter(st, lg)
	// Загружаем HTML-шаблоны из указанной директории
	r.LoadHTMLGlob("templates/*")
	// Запускаем HTTP-сервер на заданном адресе
	r.Run(cfg.ListenAddr)
	// Сохраняем текущую структуру данных в файловое хранилище
	st.Save(file)
	defer file.Close()
}
