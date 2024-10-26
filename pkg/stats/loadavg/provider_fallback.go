//go:build !linux

package loadavg

import "github.com/dima-study/monmon/pkg/stats"

const provider_platform = "fallback"

type DataProvider struct{}

func (p *DataProvider) Available() error {
	return stats.ErrNotSupported
}

func (p *DataProvider) Data() (Value, error) {
	return Value{}, stats.ErrNotSupported
}
