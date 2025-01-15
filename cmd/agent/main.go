// Команда agent запускает агент сбора метрик.
package main

import (
	"context"
	"fmt"

	"github.com/gitslim/monit/internal/agent"
	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера.
	log, err := logging.NewLogger()
	if err != nil {
		// Логгер еще недоступен поэтому panic...
		panic(fmt.Sprintf("Failed to initialize logger: %v\n", err))
	}

	// Парсинг конфига.
	cfg, err := conf.ParseConfig()
	if err != nil {
		log.Fatalf("Config parse failed: %v", err)
	}

	agent.Start(ctx, cfg, log)
}
