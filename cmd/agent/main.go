package main

import (
	"log"

	"github.com/gitslim/monit/internal/agent"
	"github.com/gitslim/monit/internal/config"
)

func main() {
	// Парсинг конфига
	conf, err := config.Parse()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	agent.Start(conf)
}
