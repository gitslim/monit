package agent

import (
	"flag"
	"time"
)

type Config struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")
	pollInterval := flag.Float64("p", 2, "Интервал сбора метрик (в секундах)")
	reportInterval := flag.Float64("r", 10, "Интервал отправки метрик на сервер (в секундах)")

	flag.Parse()

	return &Config{
		Addr:           *addr,
		PollInterval:   time.Duration(*pollInterval * float64(time.Second)),
		ReportInterval: time.Duration(*reportInterval * float64(time.Second)),
	}, nil
}
