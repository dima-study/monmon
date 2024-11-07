package providers

import (
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dima-study/monmon/internal/stats/register"
	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
	"github.com/dima-study/monmon/pkg/scheduler"
	"github.com/dima-study/monmon/pkg/stats/cpuload"
)

func init() {
	register.RegisterStat(newCPULoadProvider(), newCPULoadAggregator)
}

type cpuloadProvider struct {
	*cpuload.DataProvider
}

func newCPULoadProvider() *cpuloadProvider {
	return &cpuloadProvider{
		DataProvider: cpuload.NewDataProvider(),
	}
}

func (p *cpuloadProvider) Data() (any, error) {
	return p.DataProvider.Data()
}

const (
	fieldUserID   = "user"
	fieldUserName = "User load %"

	fieldSystemID   = "system"
	fieldSystemName = "System load %"

	fieldIdleID   = "idle"
	fieldIdleName = "Idle %"
)

func (p *cpuloadProvider) ValueToProtoRecord(val any) *v1.Record {
	v, ok := val.(cpuload.Value)
	if !ok {
		return nil
	}

	rec := v1.Record{
		Provider: p.ToProtoProvider(),
		Value: []*v1.RecordValue{
			{
				ID:    fieldUserID,
				Name:  fieldUserName,
				Value: strconv.FormatFloat(float64(v.User)/100, 'f', 2, 64),
			},
			{
				ID:    fieldSystemID,
				Name:  fieldSystemName,
				Value: strconv.FormatFloat(float64(v.System)/100, 'f', 2, 64),
			},
			{
				ID:    fieldIdleID,
				Name:  fieldIdleName,
				Value: strconv.FormatFloat(float64(v.Idle)/100, 'f', 2, 64),
			},
		},
		Time: timestamppb.New(v.T),
	}

	return &rec
}

func (p *cpuloadProvider) ToProtoProvider() *v1.Provider {
	details := v1.AvailabilityDetails{
		State:   v1.Available,
		Details: "",
	}

	errNotAvailable := register.CheckStatAvailability(p.ID())
	isDisabled, _ := register.CheckStatDisabled(p.ID())

	switch {
	case errNotAvailable != nil:
		details.State = v1.Error
		details.Details = errNotAvailable.Error()
	case isDisabled:
		details.State = v1.Disabled
	}

	protoP := v1.Provider{
		ProviderID:          p.ID(),
		ProviderName:        p.Name(),
		Platform:            p.Platform(),
		AvailabilityDetails: &details,
	}

	return &protoP
}

type cpuloadAggregator struct {
	*cpuload.CPULoad
}

func newCPULoadAggregator(n int) scheduler.Aggregator {
	return &cpuloadAggregator{
		CPULoad: cpuload.NewAggregator(n),
	}
}

func (a *cpuloadAggregator) Add(val any) {
	if v, ok := val.(cpuload.Value); ok {
		a.CPULoad.Add(v)
	}
}

func (a *cpuloadAggregator) Get(period time.Duration) (any, bool) {
	return a.CPULoad.Get(period)
}
