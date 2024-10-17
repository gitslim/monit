package main

import (
	"log"

	"github.com/gitslim/monit/internal/agent"
	"github.com/gitslim/monit/internal/agent/conf"
)

func main() {
	// Парсинг конфига
	cfg, err := conf.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	agent.Start(cfg)
}
