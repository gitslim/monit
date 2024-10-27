package conf

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   uint64 `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")
	storeInterval := flag.Uint64("i", 300, "Интервал сохранения данных на диск (в секундах)")
	fileStoragePath := flag.String("f", "/tmp/.monit/memstorage.json", "Путь до файла сохранения данных")
	restore := flag.Bool("r", true, "Флаг загрузки сохраненных данных при старте сервера")
	databaseDSN := flag.String("d", "", "Строка подключения к базе данных (DSN)")

	flag.Parse()

	cfg := &Config{
		Addr:            *addr,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		DatabaseDSN:     *databaseDSN,
	}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	// проверка конфига
	if cfg.Addr == "" {
		return nil, errors.New("адрес сервера не может быть пустым")
	}

	if cfg.FileStoragePath == "" {
		return nil, errors.New("путь до файла сохранения данных не может быть пустым")
	}

	return cfg, nil
}
