package config

const (
	defaultListenAddr    = ":8080"
	defaultStoreInterval = 300
	defaultFileStorePath = "/tmp/metrics-storage.json"
	defaultRestore       = true
)

type Config struct {
	ListenAddr      string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func New() *Config {
	return &Config{
		defaultListenAddr,
		defaultStoreInterval,
		defaultFileStorePath,
		defaultRestore,
	}
}
