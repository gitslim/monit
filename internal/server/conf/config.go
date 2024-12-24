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
	Key             string `env:"KEY"`
}

const (
	DefaultAddr            = "localhost:8080"
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "/tmp/.monit/memstorage.json"
	DefaultRestore         = true
	DefaultDatabaseDSN     = ""
	DefaultKey             = ""
)

func ParseConfig() (*Config, error) {
	addr := flag.String("a", DefaultAddr, "Адрес сервера (в формате host:port)")
	storeInterval := flag.Uint64("i", DefaultStoreInterval, "Интервал сохранения данных на диск (в секундах)")
	fileStoragePath := flag.String("f", DefaultFileStoragePath, "Путь до файла сохранения данных")
	restore := flag.Bool("r", DefaultRestore, "Флаг загрузки сохраненных данных при старте сервера")
	databaseDSN := flag.String("d", DefaultDatabaseDSN, "Строка подключения к базе данных (DSN)")
	key := flag.String("k", DefaultKey, "Ключ шифрования")

	flag.Parse()

	cfg := &Config{
		Addr:            *addr,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		DatabaseDSN:     *databaseDSN,
		Key:             *key,
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
