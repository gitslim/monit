package testhelpers

import (
	"encoding/json"
	"strings"

	"github.com/gitslim/monit/internal/entities"
)

func JSONToMetricDTO(jsonData string) ([]*entities.MetricDTO, error) {
	var metrics []*entities.MetricDTO
	err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&metrics)
	return metrics, err
}
