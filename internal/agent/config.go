package agent

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Addr           string  `env:"ADDRESS"`
	PollInterval   float64 `env:"POLL_INTERVAL"`
	ReportInterval float64 `env:"REPORT_INTERVAL"`
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")
	pollInterval := flag.Float64("p", 2, "Интервал сбора метрик (в секундах)")
	reportInterval := flag.Float64("r", 10, "Интервал отправки метрик на сервер (в секундах)")

	flag.Parse()

	cfg := &Config{
		Addr:           *addr,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
	}

	env.Parse(cfg)
	return cfg, nil
}
