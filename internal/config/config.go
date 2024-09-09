package config

import "time"

type Server struct {
	Host string
	Port uint32
}

type Config struct {
	Server         *Server
	PollInterval   time.Duration
	ReportInterval time.Duration
	Debug          bool
}

func Parse() (*Config, error) {
	return &Config{
		Server: &Server{
			Host: "localhost",
			Port: 8080},
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}, nil
}
