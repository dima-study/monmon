package start

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/dima-study/monmon/internal/server"
	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
	"github.com/dima-study/monmon/pkg/logger"
)

type (
	StartServerFunc func() error
	StopServerFunc  func(context.Context) error
)

func CreateServer(logger *logger.Logger, cfg *Config) (StartServerFunc, StopServerFunc) {
	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			server.GetStatsStreamRequestValidator(
				[2]int64{int64(cfg.Service.MinInterval), int64(cfg.Service.MaxInterval)},
				[2]int64{int64(cfg.Service.MinPeriod), int64(cfg.Service.MaxPeriod)},
			),
		),
	)

	v1.RegisterStatsServiceServer(s, server.NewStatsService(logger, Schedule))

	start := func() error {
		listenAddr := net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port)
		lstn, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}

		logger.Info("start gRPC server", "address", listenAddr)

		return s.Serve(lstn)
	}

	stop := func(ctx context.Context) error {
		done := make(chan struct{}, 1)
		go func() {
			defer close(done)
			s.GracefulStop()
		}()

		select {
		case <-ctx.Done():
			s.Stop()
			return fmt.Errorf("could not gracefully shutdown gRPC server: %w", ctx.Err())
		case <-done:
		}

		return nil
	}

	return start, stop
}
