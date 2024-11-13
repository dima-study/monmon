package loadavg

import (
	"context"
	"time"
)

const (
	providerID   = "loadavg"
	providerName = "Average system load"
)

type Value struct {
	One     int // усреднённое значение загрузки системы за 1мин (man 5 proc)
	Five    int // усреднённое значение загрузки системы за 5мин (man 5 proc)
	Fifteen int // усреднённое значение загрузки системы за 15мин (man 5 proc)

	T time.Time // когда было получено значение
}

var emptyValue = Value{}

func (v Value) IsEmpty() bool {
	return v == emptyValue
}

func (p *DataProvider) String() string {
	return providerID + "(" + providerPlatform + ")"
}

func (p *DataProvider) ID() string {
	return providerID
}

func (p *DataProvider) Name() string {
	return providerName
}

func (p *DataProvider) Platform() string {
	return providerPlatform
}

func (p *DataProvider) Cleanup(ctx context.Context) error {
	return nil
}
