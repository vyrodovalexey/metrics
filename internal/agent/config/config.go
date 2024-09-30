package config

const (
	serverAddr            = "localhost:8080"
	defaultReportInterval = 10
	defaultPoolInterval   = 2
)

type Config struct {
	EndpointAddr   string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
}

func New() *Config {
	return &Config{
		serverAddr,
		defaultReportInterval,
		defaultPoolInterval,
	}
}
