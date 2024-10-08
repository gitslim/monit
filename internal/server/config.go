package server

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   uint64 `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")
	storeInterval := flag.Uint64("i", 300, "Интервал сохранения данных на диск (в секундах)")
	fileStoragePath := flag.String("f", "/tmp/.monit/memstorage.json", "Путь до файла сохранения данных")
	restore := flag.Bool("r", true, "Флаг загрузки сохраненных данных при старте сервера")

	flag.Parse()

	cfg := &Config{
		Addr:            *addr,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
	}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
