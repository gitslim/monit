package entities

import (
	"encoding/json"
	"fmt"

	"github.com/gitslim/monit/internal/errs"
)

// MetricRequestDTO содержит данные запроса на получение метрики.
type MetricRequestDTO struct {
	ID    string     `json:"id"`   // имя метрики
	MType MetricType `json:"type"` // параметр, принимающий значение gauge или counter
}

func (s MetricRequestDTO) String() string {
	return fmt.Sprintf("MetricRequestDTO: {ID: %s, MType: %s}", s.ID, s.MType)
}

// UnmarshalJSON возвращает JSON-десериализованную MetricRequestDTO.
func (s *MetricRequestDTO) UnmarshalJSON(data []byte) error {
	var err error

	var temp struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	if err = json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.ID == "" {
		return errs.ErrMetricNotFound
	}

	switch temp.Type {
	case "counter":
		s.MType = Counter
	case "gauge":
		s.MType = Gauge
	default:
		return errs.ErrInvalidMetricType
	}

	s.ID = temp.ID

	return nil
}

func NewMetricRequestDTO(mName string, mType string) (*MetricRequestDTO, error) {
	var t MetricType
	switch mType {
	case "counter":
		t = Counter
	case "gauge":
		t = Gauge
	default:
		return nil, errs.ErrInvalidMetricType
	}

	return &MetricRequestDTO{
		ID:    mName,
		MType: t,
	}, nil
}
