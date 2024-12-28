package conf

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Значения по умолчанию для конфигурации
const (
	DefaultAddr           = "localhost:8080"
	DefaultPollInterval   = 2
	DefaultReportInterval = 10
	DefaultKey            = ""
	DefaultRateLimit      = 10
)

// Config представляет конфигурацию агента сбора метрик
type Config struct {
	Addr           string `env:"ADDRESS"`
	PollInterval   uint64 `env:"POLL_INTERVAL"`
	ReportInterval uint64 `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      uint64 `env:"RATE_LIMIT"`
}

// ParseConfig парсит конфигурацию из флагов и переменных окружения
func ParseConfig() (*Config, error) {
	addr := flag.String("a", DefaultAddr, "Адрес сервера (в формате host:port)")
	pollInterval := flag.Uint64("p", DefaultPollInterval, "Интервал сбора метрик (в секундах)")
	reportInterval := flag.Uint64("r", DefaultReportInterval, "Интервал отправки метрик на сервер (в секундах)")
	key := flag.String("k", DefaultKey, "Ключ шифрования")
	rateLimit := flag.Uint64("l", DefaultRateLimit, "Лимит одновременно исходящих запросов на отправку метрик")

	flag.Parse()

	cfg := &Config{
		Addr:           *addr,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		Key:            *key,
		RateLimit:      *rateLimit,
	}

	// Парсинг конфига
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	// проверка конфига
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// valiadteConfig - проверка конфига на корректность
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
