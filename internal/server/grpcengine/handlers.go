package grpcengine

import (
	"context"

	"github.com/gitslim/monit/internal/entities"
	pb "github.com/gitslim/monit/pkg/pb/v1"
)

func metricToPB(metric entities.Metric) *pb.Metric {
	delta, ok := metric.GetValue().(int64)
	if !ok {
		delta = 0
	}

	value, ok := metric.GetValue().(float64)
	if !ok {
		value = 0
	}

	return &pb.Metric{
		Id:    metric.GetName(),
		Type:  metric.GetType().String(),
		Delta: delta,
		Value: value,
	}
}

func (s *MonitServer) ListMetrics(ctx context.Context, in *pb.ListMetricsRequest) (*pb.ListMetricsResponse, error) {
	metrics := []*pb.Metric{}

	all, err := s.svc.GetAllMetrics()
	if err != nil {
		return nil, err
	}

	for _, metric := range all {
		metrics = append(metrics, metricToPB(metric))
	}
	var res = &pb.ListMetricsResponse{Metrics: metrics}
	return res, nil
}

func (s *MonitServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	metric, err := s.svc.GetMetric(in.Id, in.Type)
	if err != nil {
		return nil, err
	}
	var res = &pb.GetMetricResponse{Metric: metricToPB(metric)}
	return res, nil
}

func (s *MonitServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	dto := entities.MetricDTO{
		ID:    in.Metric.Id,
		MType: in.Metric.Type,
		Delta: in.Metric.Delta,
		Value: in.Metric.Value,
	}
	err := s.svc.UpdateMetric(dto)
	if err != nil {
		return nil, err
	}
	var res = &pb.UpdateMetricResponse{Message: "success"}
	return res, nil
}

func (s *MonitServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	dtos := []*entities.MetricDTO{}
	for _, metric := range in.Metrics {
		dtos = append(dtos, &entities.MetricDTO{
			ID:    metric.Id,
			MType: metric.Type,
			Delta: metric.Delta,
			Value: metric.Value,
		})
	}

	err := s.svc.BatchUpdateMetrics(dtos)
	if err != nil {
		return nil, err
	}
	var res = &pb.UpdateMetricsResponse{Message: "success"}
	return res, nil
}

func (s *MonitServer) PingStorage(ctx context.Context, in *pb.PingStorageRequest) (*pb.PingStorageResponse, error) {
	err := s.svc.PingStorage()
	if err != nil {
		return nil, err
	}
	var res = &pb.PingStorageResponse{}
	return res, nil
}
