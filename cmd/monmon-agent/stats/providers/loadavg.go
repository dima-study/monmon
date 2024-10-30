package providers

import (
	"time"

	"github.com/dima-study/monmon/cmd/monmon-agent/stats/register"
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
