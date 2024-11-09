package cpuload

import (
	"context"
	"time"
)

const (
	providerID   = "cpuload"
	providerName = "CPU load"
)

type Value struct {
	User   int // % времени отданного пользователям
	System int // % времени отданного системе
	Idle   int // % времени бездействия

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
