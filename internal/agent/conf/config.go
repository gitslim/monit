// Package conf создает конфигурацию агента сбора метрик, используя флаги, переменные окружения и значения по умолчанию.
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
	DefaultAddr           = "localhost:8080"
	DefaultPollInterval   = 2
	DefaultReportInterval = 10
	DefaultKey            = ""
	DefaultRateLimit      = 10
	DefaultCryptoKey      = ""
	DefaultConfig         = ""
)

// Config представляет конфигурацию агента сбора метрик.
type Config struct {
	Addr           string `env:"ADDRESS" json:"address"`
	PollInterval   uint64 `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval uint64 `env:"REPORT_INTERVAL" json:"report_interval"`
	Key            string `env:"KEY" json:"key"`
	RateLimit      uint64 `env:"RATE_LIMIT" json:"rate_limit"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath     string `env:"CONFIG" json:"-"`
}

// ParseConfig парсит конфигурацию из json-конфига, флагов и переменных окружения.
func ParseConfig() (*Config, error) {
	configPath := flag.String("config", DefaultConfig, "Путь до конфигурационного файла (JSON)")
	addr := flag.String("a", DefaultAddr, "Адрес сервера (host:port)")
	pollInterval := flag.Uint64("p", DefaultPollInterval, "Интервал сбора метрик (сек)")
	reportInterval := flag.Uint64("r", DefaultReportInterval, "Интервал отправки метрик (сек)")
	key := flag.String("k", DefaultKey, "Ключ шифрования")
	rateLimit := flag.Uint64("l", DefaultRateLimit, "Лимит запросов")
	cryptoKey := flag.String("crypto-key", DefaultCryptoKey, "Публичный ключ шифрования")

	// Парсим флаги
	flag.Parse()

	// Загружаем конфиг из JSON если путь указан
	cfg := Config{
		Addr:           DefaultAddr,
		PollInterval:   DefaultPollInterval,
		ReportInterval: DefaultReportInterval,
		Key:            DefaultKey,
		RateLimit:      DefaultRateLimit,
		CryptoKey:      DefaultCryptoKey,
		ConfigPath:     *configPath,
	}

	if *configPath != "" {
		if err := loadConfigFromJSON(*configPath, &cfg); err != nil {
			return nil, err
		}
	}

	// Перезаписываем значениями из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга env: %w", err)
	}

	// Перезаписываем значениями флагов (если они были переданы)
	if flag.Lookup("a").Value.String() != DefaultAddr {
		cfg.Addr = *addr
	}
	if flag.Lookup("p").Value.String() != fmt.Sprint(DefaultPollInterval) {
		cfg.PollInterval = *pollInterval
	}
	if flag.Lookup("r").Value.String() != fmt.Sprint(DefaultReportInterval) {
		cfg.ReportInterval = *reportInterval
	}
	if flag.Lookup("k").Value.String() != DefaultKey {
		cfg.Key = *key
	}
	if flag.Lookup("l").Value.String() != fmt.Sprint(DefaultRateLimit) {
		cfg.RateLimit = *rateLimit
	}
	if flag.Lookup("crypto-key").Value.String() != DefaultCryptoKey {
		cfg.CryptoKey = *cryptoKey
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
	if cfg.Addr == "" {
		return errors.New("адрес сервера не может быть пустым")
	}

	if cfg.PollInterval == 0 {
		return errors.New("интервал сбора метрик не может быть равен 0")
	}

	if cfg.ReportInterval == 0 {
		return errors.New("интервал отправки метрик на сервер не может быть равен 0")
	}

	if cfg.RateLimit == 0 {
		return errors.New("лимит одновременно исходящих запросов на отправку метрик не может быть равен 0")
	}

	return nil
}
