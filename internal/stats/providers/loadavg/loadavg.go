package providers

import (
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dima-study/monmon/internal/stats/register"
	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
	"github.com/dima-study/monmon/pkg/scheduler"
	"github.com/dima-study/monmon/pkg/stats/loadavg"
)

func init() {
	register.RegisterStat(newLoadavgProvider(), newLoadavgAggregator)
}

type loadavgProvider struct {
	*loadavg.DataProvider
}

func newLoadavgProvider() *loadavgProvider {
	return &loadavgProvider{
		DataProvider: &loadavg.DataProvider{},
	}
}

func (p *loadavgProvider) Data() (any, error) {
	return p.DataProvider.Data()
}

const (
	field1MinID   = "1min"
	field1MinName = "1 minute"

	field5MinID   = "5min"
	field5MinName = "5 minutes"

	field15MinID   = "15min"
	field15MinName = "15 minutes"
)

func (p *loadavgProvider) ValueToProtoRecord(val any) *v1.Record {
	v, ok := val.(loadavg.Value)
	if !ok {
		return nil
	}

	rec := v1.Record{
		Provider: p.ToProtoProvider(),
		Value: []*v1.RecordValue{
			{
				ID:    field1MinID,
				Name:  field1MinName,
				Value: strconv.FormatFloat(float64(v.One)/100, 'f', 2, 64),
			},
			{
				ID:    field5MinID,
				Name:  field5MinName,
				Value: strconv.FormatFloat(float64(v.Five)/100, 'f', 2, 64),
			},
			{
				ID:    field15MinID,
				Name:  field15MinName,
				Value: strconv.FormatFloat(float64(v.Fifteen)/100, 'f', 2, 64),
			},
		},
		Time: timestamppb.New(v.T),
	}

	return &rec
}

func (p *loadavgProvider) ToProtoProvider() *v1.Provider {
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

type loadavgAggregator struct {
	*loadavg.LoadAvg
}

func newLoadavgAggregator(n int) scheduler.Aggregator {
	return &loadavgAggregator{
		LoadAvg: loadavg.NewAggregator(n),
	}
}

func (a *loadavgAggregator) Add(val any) {
	if v, ok := val.(loadavg.Value); ok {
		a.LoadAvg.Add(v)
	}
}

func (a *loadavgAggregator) Get(period time.Duration) (any, bool) {
	return a.LoadAvg.Get(period)
}
