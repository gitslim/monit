package main

import (
	"log"

	"github.com/gitslim/monit/internal/agent"
)

func main() {
	// Парсинг конфига
	cfg, err := agent.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	agent.Start(cfg)
}
