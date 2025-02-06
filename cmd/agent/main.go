// Команда agent запускает агент сбора метрик.
package main

import (
	"fmt"

	"github.com/gitslim/monit/internal/agent"
	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/logging"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", getOrDefault(buildVersion))
	fmt.Printf("Build date: %s\n", getOrDefault(buildDate))
	fmt.Printf("Build commit: %s\n", getOrDefault(buildCommit))
}

func getOrDefault(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func main() {
	// вывод информации о билде.
	printBuildInfo()

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

	agent.Start(cfg, log)
}
