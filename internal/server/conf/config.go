// Package conf создает конфигурацию сервера метрик, используя флаги, переменные окружения и значения по умолчанию.
package conf

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

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
	DefaultConfig          = ""
	DefaultTrustedSubnet   = ""
)

// Config представляет конфигурацию сервера.
type Config struct {
	Addr            string `env:"ADDRESS" json:"address"`
	StoreInterval   uint64 `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	Key             string `env:"KEY" json:"key"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath      string `env:"CONFIG" json:"-"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// ParseConfig парсит конфигурацию из флагов и переменных окружения.
func ParseConfig() (*Config, error) {
	configPath := flag.String("config", DefaultConfig, "Путь до конфигурационного файла (JSON)")
	addr := flag.String("a", DefaultAddr, "Адрес сервера (в формате host:port)")
	storeInterval := flag.Uint64("i", DefaultStoreInterval, "Интервал сохранения данных на диск (в секундах)")
	fileStoragePath := flag.String("f", DefaultFileStoragePath, "Путь до файла сохранения данных")
	restore := flag.Bool("r", DefaultRestore, "Флаг загрузки сохраненных данных при старте сервера")
	databaseDSN := flag.String("d", DefaultDatabaseDSN, "Строка подключения к базе данных (DSN)")
	key := flag.String("k", DefaultKey, "Ключ шифрования")
	cryptoKey := flag.String("crypto-key", DefaultCryptoKey, "Приватный ключ шифрования")
	trustedSubnet := flag.String("t", DefaultTrustedSubnet, "Адрес доверенной подсети в формате CIDR")

	// Парсим флаги
	flag.Parse()

	// Загружаем конфиг из JSON если путь указан
	cfg := Config{
		Addr:            DefaultAddr,
		StoreInterval:   DefaultStoreInterval,
		FileStoragePath: DefaultFileStoragePath,
		Restore:         DefaultRestore,
		DatabaseDSN:     DefaultDatabaseDSN,
		Key:             DefaultKey,
		CryptoKey:       DefaultCryptoKey,
		TrustedSubnet:   DefaultTrustedSubnet,
		ConfigPath:      *configPath,
	}

	if *configPath != "" {
		if err := loadConfigFromJSON(*configPath, &cfg); err != nil {
			return nil, err
		}
	}

	// Перезаписываем значениями флагов (если они были переданы)
	if flag.Lookup("a").Value.String() != DefaultAddr {
		cfg.Addr = *addr
	}
	if flag.Lookup("i").Value.String() != fmt.Sprint(DefaultStoreInterval) {
		cfg.StoreInterval = *storeInterval
	}
	if flag.Lookup("f").Value.String() != fmt.Sprint(DefaultFileStoragePath) {
		cfg.FileStoragePath = *fileStoragePath
	}
	if flag.Lookup("r").Value.String() != fmt.Sprint(DefaultRestore) {
		cfg.Restore = *restore
	}
	if flag.Lookup("d").Value.String() != DefaultDatabaseDSN {
		cfg.DatabaseDSN = *databaseDSN
	}
	if flag.Lookup("k").Value.String() != DefaultKey {
		cfg.Key = *key
	}
	if flag.Lookup("crypto-key").Value.String() != DefaultCryptoKey {
		cfg.CryptoKey = *cryptoKey
	}
	if flag.Lookup("t").Value.String() != DefaultTrustedSubnet {
		cfg.TrustedSubnet = *trustedSubnet
	}

	// Перезаписываем значениями из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга env: %w", err)
	}

	// Валидация
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// loadConfigFromJSON загружает конфигурацию из JSON-файла.
func loadConfigFromJSON(path string, cfg *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	defer func() { _ = file.Close() }()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %w", err)
	}
	return nil
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
