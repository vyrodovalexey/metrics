package main

import (
	"github.com/vyrodovalexey/metrics/internal/server/config"
	"github.com/vyrodovalexey/metrics/internal/server/logging"
	"github.com/vyrodovalexey/metrics/internal/server/memstorage"
	"github.com/vyrodovalexey/metrics/internal/server/routing"
	storage2 "github.com/vyrodovalexey/metrics/internal/server/storage"
	"go.uber.org/zap"
	"os"
)

const (
	dbConnectionTimout = 5
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
		"Database connection string", cfg.DatabaseDSN,
	)

	//	ctx := context.Background()
	//	conn, err := pgx.Connect(ctx, cfg.DatabaseDSN)
	//	//	return conn, err
	//	if err != nil {
	//		lg.Panicw("Initializing database connection...",
	//			"Database connection error:", err,
	//		)
	//	}

	// Инициализируем интерфейс и структуру хранения данных
	var st storage2.Storage = &memstorage.MemStorageWithAttributes{}
	file, errr := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if errr != nil {
		// Логируем ошибку, если открытие/создание файла не удалось
		lg.Panicw("Initializing file storage...",
			"Error creating file:", errr,
		)
		return
	}
	lg.Infow("File storage initialized")
	defer file.Close()

	// Проверяем, нужно ли загружать файл хранилища
	// Если нет, инициализируем новое
	if cfg.Restore {
		err := st.LoadMemStorage(cfg.FileStoragePath, cfg.StoreInterval)
		if err != nil {
			// Логируем ошибку, если открытие/создание файла не удалось
			lg.Panicw("Initializing file storage...",
				"Error load data from file:", err,
			)
			return
		}
	} else {
		// Логируем ошибку, если открытие/создание файла не удалось
		err := st.NewMemStorage(cfg.FileStoragePath, cfg.StoreInterval)
		if err != nil {
			lg.Panicw("Initializing file storage...",
				"Error creating file:", err,
			)
			return
		}
	}
	lg.Infow("File storage initialized")

	// При ассинхронном режиме запускаем фоновое сохранение структуры данных
	// и назначение переменной owf
	if cfg.StoreInterval > 0 {
		go st.SaveAsync()
	}

	// Инициализируем маршрутизатор с хранилищем и логированием
	r := routing.SetupRouter(lg)
	// Настраиваем маршрутизацию
	routing.ConfigureRouting(r, st)
	// Загружаем HTML-шаблоны из указанной директории
	r.LoadHTMLGlob("templates/*")
	// Запускаем HTTP-сервер на заданном адресе
	r.Run(cfg.ListenAddr)
	// Сохраняем текущую структуру данных в файловое хранилище
	err := st.Save()
	if err != nil {
		// Логируем ошибку, если открытие/создание файла не удалось
		lg.Panicw("Saving data...",
			"Error saving data:", err,
		)
		return
	}
	st.Close()
}
