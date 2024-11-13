package server

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
)

func GetStatsStreamRequestValidator(interval [2]int64, period [2]int64) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &getStatsStreamRequestValidator{
			ServerStream: ss,
			info:         info,
			interval:     interval,
			period:       period,
		})
	}
}

type getStatsStreamRequestValidator struct {
	grpc.ServerStream

	info *grpc.StreamServerInfo

	interval [2]int64 // [min, max]
	period   [2]int64 // [min, max]
}

func (v *getStatsStreamRequestValidator) RecvMsg(m any) error {
	if err := v.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	if err := v.validate(m); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func (v *getStatsStreamRequestValidator) validate(m any) error {
	r, ok := m.(*v1.GetStatsRequest)
	if !ok {
		return nil
	}

	if r.Interval < v.interval[0] || v.interval[1] < r.Interval {
		return fmt.Errorf(
			"invalid param 'interval': got %d, must be in [%d, %d]",
			r.Interval, v.interval[0], v.interval[1],
		)
	}

	if r.Period < v.period[0] || v.period[1] < r.Period {
		return fmt.Errorf("invalid param 'period': got %d, must be in [%d, %d]",
			r.Period, v.period[0], v.period[1],
		)
	}

	return nil
}
