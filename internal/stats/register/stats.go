// register - реестр для агрегаторов и провайдеров статистики.
package register

import (
	"slices"

	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
	"github.com/dima-study/monmon/pkg/scheduler"
	"github.com/dima-study/monmon/pkg/stats"
)

// DataProvider - провайдер статистики.
//
// Добавляет в провайдер планировщика требования нескольких методов, которые используются в gRPC сервере.
type DataProvider interface {
	scheduler.DataProvider

	// ID - идентификатор продайдера
	ID() string

	// Name - название продайдера
	Name() string

	// ValueToProtoRecord преобразует данные val от провайдера в соответствующей тип gRPC message.
	// Возвращает nil, если данные val не являются данными от провайдера.
	ValueToProtoRecord(val any) *v1.Record

	// ToProtoProvider преобразует информацию о провайдере в соответствующей тип gRPC message.
	ToProtoProvider() *v1.Provider
}

// AggregatorMaker - функция/обёртка для создания агрегатора планировщика.
type AggregatorMaker func(n int) scheduler.Aggregator

type supportedStat struct {
	provider        DataProvider
	aggregatorMaker AggregatorMaker
	available       bool
	reason          error
	disabled        bool
}

var supportedStats = make(map[string]supportedStat)

// RegisterStat регистрирует доступный провайдер.
func RegisterStat(provider DataProvider, am AggregatorMaker) { //nolint:revive
	reason := provider.Available()

	supportedStats[provider.ID()] = supportedStat{
		provider:        provider,
		aggregatorMaker: am,
		available:       reason == nil,
		reason:          reason,
		disabled:        false,
	}
}

// CheckStatSupport возвращает ошибку, если провайдер с соответствующим ID не поддерживается
// (т.е. не был зарегистрирован).
func CheckStatSupport(providerID string) error {
	_, exists := supportedStats[providerID]
	if !exists {
		return stats.ErrNotSupported
	}

	return nil
}

// CheckStatAvailability возвращает ошибку, если поддерживаемый провайдер с соответствующим ID не доступен.
// Также будет возвращена ошибка от вызова CheckStatSupport для проверки поддерживаемого провайдера.
func CheckStatAvailability(providerID string) error {
	if err := CheckStatSupport(providerID); err != nil {
		return err
	}

	return supportedStats[providerID].reason
}

// CheckStatDisabled возвращает информацию о том отключён ли поддерживаемый провайдер с соответствующим ID.
// Также будет возвращена ошибка от вызова CheckStatSupport для проверки поддерживаемого провайдера.
func CheckStatDisabled(providerID string) (bool, error) {
	if err := CheckStatSupport(providerID); err != nil {
		return true, err
	}

	return supportedStats[providerID].disabled, nil
}

// DisableStat отмечает поддерживаемый провайдер с соответствующим ID как отключённый.
// Также будет возвращена ошибка от вызова CheckStatSupport для проверки поддерживаемого провайдера.
func DisableStat(providerID string) error {
	if err := CheckStatSupport(providerID); err != nil {
		return err
	}

	p := supportedStats[providerID]
	p.disabled = true
	supportedStats[providerID] = p

	return nil
}

// SupportedStats возвращает список (отсортированный) с идентификаторами поддерживаемых провайдеров.
func SupportedStats() []string {
	names := make([]string, 0, len(supportedStats))
	for k := range supportedStats {
		names = append(names, k)
	}

	slices.Sort(names)
	return names
}

// GetProvider возвращает соответствующий поддерживаемый провайдер или ошибку, если провайдер не поддерживается.
func GetProvider(providerID string) (DataProvider, error) {
	if err := CheckStatSupport(providerID); err != nil {
		return nil, err
	}

	return supportedStats[providerID].provider, nil
}

// GetAggregatorMaker возвращает соответствующий AggregatorMaker для поддерживаемого провайдера
// или ошибку, если провайдер не поддерживается.
func GetAggregatorMaker(providerID string) (AggregatorMaker, error) {
	if err := CheckStatSupport(providerID); err != nil {
		return nil, err
	}

	return supportedStats[providerID].aggregatorMaker, nil
}
