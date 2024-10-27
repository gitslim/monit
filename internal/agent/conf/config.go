package conf

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	PollInterval   uint64 `env:"POLL_INTERVAL"`
	ReportInterval uint64 `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")
	pollInterval := flag.Uint64("p", 2, "Интервал сбора метрик (в секундах)")
	reportInterval := flag.Uint64("r", 10, "Интервал отправки метрик на сервер (в секундах)")
	key := flag.String("k", "", "Ключ шифрования")

	flag.Parse()

	cfg := &Config{
		Addr:           *addr,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		Key:            *key,
	}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	// проверка конфига
	if cfg.Addr == "" {
		return nil, errors.New("адрес сервера не может быть пустым")
	}

	if cfg.PollInterval == 0 {
		return nil, errors.New("интервал сбора метрик не может быть равен 0")
	}

	if cfg.ReportInterval == 0 {
		return nil, errors.New("интервал отправки метрик на сервер не может быть равен 0")
	}

	return cfg, nil
}
