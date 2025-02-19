package grpcengine

import (
	"context"
	"net"
	"sync"

	"github.com/gitslim/monit/internal/logging"
	"github.com/gitslim/monit/internal/server/conf"
	"github.com/gitslim/monit/internal/services"
	pb "github.com/gitslim/monit/pkg/pb/v1"
	"google.golang.org/grpc"
)

type MonitServer struct {
	pb.UnimplementedMonitServer
	cfg *conf.Config
	log *logging.Logger
}

func StartGRPCServer(ctx context.Context, cfg *conf.Config, log *logging.Logger, wg *sync.WaitGroup, metricService *services.MetricService) {
	defer wg.Done()

	// Создаем сервис gRPC
	s := &MonitServer{cfg: cfg, log: log}

	// Создаем TCP-сокет
	lis, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		s.log.Fatalf("gRPC server failed to listen: %v", err)
	}

	// Создаем сервер gRPC
	grpcServer := grpc.NewServer()
	// Регистрируем сервис
	pb.RegisterMonitServer(grpcServer, s)

	// Запускаем сервер gRPC
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			s.log.Fatalf("gRPC server failed to serve: %v", err)
		}
	}()

	s.log.Infof("gRPC server is running on %s", cfg.Addr)

	// Ожидаем сигнал завершения.
	<-ctx.Done()

	// Останавливаем сервер gRPC
	grpcServer.Stop()
	s.log.Info("gRPC server stopped.")
}
