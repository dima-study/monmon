package loadavg

import "time"

const provider_id = "loadavg"

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
	return provider_id + "(" + provider_platform + ")"
}

func (p *DataProvider) ID() string {
	return provider_id
}

func (p *DataProvider) Platform() string {
	return provider_platform
}
