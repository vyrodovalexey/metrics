package config

const (
	defaultListenAddr      = ":8080"                     // Адрес для прослушивания по умолчанию (:8080)
	defaultStoreInterval   = 300                         // Интервал сохранения данных по умолчанию (в секундах)
	defaultFileStorePath   = "/tmp/metrics-storage.json" // Путь к файлу хранения метрик по умолчанию
	defaultRestore         = true                        // Флаг загрузки файла данных при запуске по умолчанию (включено)
	defaultDatabaseDSN     = ""                          // Строка подключения к базе данных по умолчанию (пустая строка)
	defaultDatabaseTimeout = 0                           // Таймаут подключения к базе данных по умолчанию (в секундах)
)

// Config Структура для хранения конфигурации
type Config struct {
	ListenAddr      string `env:"ADDRESS"`
	StoreInterval   uint   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	DatabaseTimeout uint   `env:"DATABASE_TIMEOUT"`
}

// New Функция для создания нового экземпляра конфигурации
func New() *Config {
	return &Config{
		defaultListenAddr,
		defaultStoreInterval,
		defaultFileStorePath,
		defaultRestore,
		defaultDatabaseDSN,
		defaultDatabaseTimeout,
	}
}
