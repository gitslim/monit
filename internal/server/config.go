package server

import (
	"flag"
)

type Config struct {
	Addr string
}

func ParseConfig() (*Config, error) {
	addr := flag.String("a", "localhost:8080", "Адрес сервера (в формате host:port)")

	flag.Parse()

	return &Config{
		Addr: *addr,
	}, nil
}
