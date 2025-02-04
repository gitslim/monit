// Package conf создает конфигурацию сервера метрик, используя флаги, переменные окружения и значения по умолчанию.
package conf

import (
	"errors"
	"flag"
	"fmt"

	env "github.com/caarlos0/env/v6"
)

// Значения по умолчанию для конфигурации.
const (
	DefaultAddr            = "localhost:8080"
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "/tmp/.monit/memstorage.json"
	DefaultRestore         = true
	DefaultDatabaseDSN     = ""
	DefaultKey             = ""
	DefaultCryptoKey       = ""
)

// Config представляет конфигурацию сервера.
type Config struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   uint64 `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
}

// ParseConfig парсит конфигурацию из флагов и переменных окружения.
func ParseConfig() (*Config, error) {
	addr := flag.String("a", DefaultAddr, "Адрес сервера (в формате host:port)")
	storeInterval := flag.Uint64("i", DefaultStoreInterval, "Интервал сохранения данных на диск (в секундах)")
	fileStoragePath := flag.String("f", DefaultFileStoragePath, "Путь до файла сохранения данных")
	restore := flag.Bool("r", DefaultRestore, "Флаг загрузки сохраненных данных при старте сервера")
	databaseDSN := flag.String("d", DefaultDatabaseDSN, "Строка подключения к базе данных (DSN)")
	key := flag.String("k", DefaultKey, "Ключ шифрования")
	cryptoKey := flag.String("crypto-key", DefaultCryptoKey, "Приватный ключ шифрования")

	flag.Parse()

	cfg := &Config{
		Addr:            *addr,
		StoreInterval:   *storeInterval,
		FileStoragePath: *fileStoragePath,
		Restore:         *restore,
		DatabaseDSN:     *databaseDSN,
		Key:             *key,
		CryptoKey:       *cryptoKey,
	}

	// Парсинг конфига.
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	// Проверка конфига.
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// valiadteConfig - проверка конфига на корректность.
func validateConfig(cfg *Config) error {
	// Проверка конфига.
	if cfg.Addr == "" {
		return errors.New("адрес сервера не может быть пустым")
	}

	if cfg.FileStoragePath == "" {
		return errors.New("путь до файла сохранения данных не может быть пустым")
	}

	return nil
}
