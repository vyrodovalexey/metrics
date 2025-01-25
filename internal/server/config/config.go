package config

const (
	defaultListenAddr    = ":8080"                // Адрес для прослушивания по умолчанию (:8080)
	defaultStoreInterval = 300                    // Интервал сохранения данных по умолчанию (в секундах)
	defaultFileStorePath = "metrics-storage.json" // Путь к файлу хранения метрик по умолчанию
	defaultRestore       = true                   // Флаг загрузки файла данных при запуске по умолчанию (включено)
)

// Config Структура для хранения конфигурации
type Config struct {
	ListenAddr      string `env:"ADDRESS"`
	StoreInterval   uint   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// New Функция для создания нового экземпляра конфигурации
func New() *Config {
	return &Config{
		defaultListenAddr,
		defaultStoreInterval,
		defaultFileStorePath,
		defaultRestore,
		"",
	}
}
