package config

const (
	serverAddr            = "localhost:8080" // Адрес сервера по умолчанию
	defaultReportInterval = 10               // Интервал отправки метрик по умолчанию (в секундах)
	defaultPoolInterval   = 2                // Интервал опроса метрик по умолчанию (в секундах)
	defaultBatchMode      = true
)

// Config Структура для хранения конфигурации
type Config struct {
	EndpointAddr   string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
	BatchMode      bool   `env:"BATCH_MODE"`
}

// New Функция для создания нового экземпляра конфигурации
func New() *Config {
	return &Config{
		serverAddr,
		defaultReportInterval,
		defaultPoolInterval,
		defaultBatchMode,
	}
}
