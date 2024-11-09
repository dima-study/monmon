package server

import (
	"context"
	"time"

	"github.com/dima-study/monmon/internal/stats/register"
	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
	"github.com/dima-study/monmon/pkg/logger"
)

type ScheduleFunc func(ctx context.Context, interval time.Duration, period time.Duration) <-chan *v1.Record

type StatsService struct {
	logger     *logger.Logger
	scheduleFn ScheduleFunc

	v1.UnimplementedStatsServiceServer
}

func NewStatsService(logger *logger.Logger, schedule ScheduleFunc) *StatsService {
	return &StatsService{
		logger:     logger,
		scheduleFn: schedule,
	}
}

func (srv *StatsService) GetStats(req *v1.GetStatsRequest, stream v1.StatsService_GetStatsServer) error {
	ch := srv.scheduleFn(
		stream.Context(),
		time.Duration(req.Interval)*time.Second,
		time.Duration(req.Period)*time.Second,
	)

	for rec := range ch {
		// Т.к. ch буферезированный, то в нём могут быть данные даже после завершения сеанса.
		select {
		case <-stream.Context().Done():
			return nil
		default:
		}

		m := &v1.GetStatsResponse{
			Record: rec,
		}

		if err := stream.SendMsg(m); err != nil {
			srv.logger.Error("can't send stat", "error", err)
			return nil
		}
	}

	return nil
}

func (srv *StatsService) GetSupportedStats(
	context.Context,
	*v1.GetSupportedStatsRequest,
) (*v1.GetSupportedStatsResponse, error) {
	providers := []*v1.Provider{}

	for _, providerID := range register.SupportedStats() {
		p, err := register.GetProvider(providerID)
		if err != nil {
			srv.logger.Error(
				"can't get provider",
				"providerID", providerID,
				"error", err,
			)
			continue
		}

		providers = append(providers, p.ToProtoProvider())
	}

	return &v1.GetSupportedStatsResponse{Providers: providers}, nil
}
