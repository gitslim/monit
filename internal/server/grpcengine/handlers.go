package grpcengine

import (
	"context"

	pb "github.com/gitslim/monit/pkg/pb/v1"
)

func (s *MonitServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var res = &pb.UpdateMetricResponse{}
	return res, nil
}
