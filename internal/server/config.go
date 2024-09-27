package server

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Addr string `env:"ADDRESS"`
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")

	flag.Parse()

	cfg := &Config{
		Addr: *addr,
	}

	env.Parse(cfg)
	return cfg, nil
}
