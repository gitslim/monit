package sender

import (
	"context"

	"github.com/gitslim/monit/internal/agent/conf"
	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/retry"
	pb "github.com/gitslim/monit/pkg/pb/v1"
)

func GRPCSendMetrics(ctx context.Context, cfg *conf.Config, monitClient pb.MonitClient, dtos []*entities.MetricDTO, batch bool) error {
	return retry.Retry(func() error {
		var metrics []*pb.Metric

		for _, m := range dtos {
			metric := &pb.Metric{
				Id:    m.ID,
				Type:  m.MType,
				Delta: m.Delta,
				Value: m.Value,
			}

			if batch {
				metrics = append(metrics, metric)
			} else {
				req := &pb.UpdateMetricRequest{Metric: metric}
				if _, err := monitClient.UpdateMetric(ctx, req); err != nil {
					return err
				}
			}
		}

		if batch {
			req := &pb.UpdateMetricsRequest{Metrics: metrics}
			if _, err := monitClient.UpdateMetrics(ctx, req); err != nil {
				return err
			}
		}

		return nil
	}, 3)
}
