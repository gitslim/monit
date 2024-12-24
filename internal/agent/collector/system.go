package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/gitslim/monit/internal/agent/worker"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/logging"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// CollectSystemMetrics собирает метрики системной информации и отправляет их в канал wp.Metrics.
func CollectSystemMetrics(ctx context.Context, log *logging.Logger, wp *worker.WorkerPool) {
	pollTicker := time.NewTicker(time.Duration(wp.Cfg.PollInterval * uint64(time.Second)))
	var metric *entities.MetricDTO

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			vMem, err := mem.VirtualMemory()
			if err != nil {
				log.Errorf("failed to get memory info: %v", err)
			} else {
				metric, err = entities.NewMetricDTO("TotalMemory", "gauge", float64(vMem.Total))
				if err != nil {
					log.Error("Failed to create gauge DTO: TotalMemory")
				} else {
					wp.Metrics <- *metric
				}

				metric, err = entities.NewMetricDTO("FreeMemory", "gauge", float64(vMem.Free))
				if err != nil {
					log.Error("Failed to create gauge DTO: FreeMemory")
				} else {
					wp.Metrics <- *metric
				}
			}

			cpuPercents, err := cpu.Percent(0, true)
			if err != nil {
				log.Errorf("failed to get CPU info: %v", err)
			} else {
				for i, cpuPercent := range cpuPercents {
					metricName := fmt.Sprintf("CPUutilization%d", i+1)
					metric, err = entities.NewMetricDTO(metricName, "gauge", cpuPercent)
					if err != nil {
						log.Errorf("Failed to create gauge DTO: %v", metricName)
					} else {
						wp.Metrics <- *metric
					}

				}
			}
		}
	}
}
