package start

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/dima-study/monmon/internal/stats/providers" // автоматическая регистрация провайдеров
	"github.com/dima-study/monmon/internal/stats/register"
	v1 "github.com/dima-study/monmon/pkg/api/proto/stats/v1"
	"github.com/dima-study/monmon/pkg/logger"
	"github.com/dima-study/monmon/pkg/scheduler"
)

// Coordinator содержит в себе информацию о координаторе планировщика, агрегаторе
// и функцию конвертации значения от агрегатора в Record для дальнейшей передачи клиенту.
type Coordinator struct {
	c                    *scheduler.Coordinator
	logger               *logger.Logger
	agg                  scheduler.Aggregator
	valueToProtoRecordFn func(val any) *v1.Record
}

var coordinators = []Coordinator{}

// InitCoordinators инициализирует поддерживаемые и доступные для использования координаторы планировщика.
// Провайдеры будут запущены с точностью accuracy.
func InitCoordinators(ctx context.Context, logger *logger.Logger, accuracy int) error {
	statsList := register.SupportedStats()

	for _, providerID := range statsList {
		logger.Info("registering coordinator", "provider", providerID)

		if err := register.CheckStatAvailability(providerID); err != nil {
			logger.Info("provider is not available", "provider", providerID, "reason", err.Error())
			continue
		}

		if disabled, _ := register.CheckStatDisabled(providerID); disabled {
			logger.Info("provider is disabled", "provider", providerID)
			continue
		}

		provider, _ := register.GetProvider(providerID)
		aggMaker, _ := register.GetAggregatorMaker(providerID)

		crd, err := initCoordinator(
			ctx,
			logger.With("coordinator", providerID),
			provider,
			accuracy,
			aggMaker,
		)
		if err != nil {
			return fmt.Errorf("failed to init coordinator for provider '%s': %w", providerID, err)
		}

		coordinators = append(coordinators, crd)
	}

	return nil
}

func ResetCoordinators(ctx context.Context) error {
	errCh := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(coordinators))

	for i := range len(coordinators) {
		crd := coordinators[i]
		go func() {
			defer wg.Done()

			if err := crd.c.Reset(ctx); err != nil {
				errCh <- err
			}
		}()
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	var err error
	for e := range errCh {
		err = errors.Join(err, e)
	}

	return err
}

// Grower инерфейс для увеличения размера агрегатора.
type Grower interface {
	// Grow увеличивает размер агрегатора до n
	Grow(n int)
}

// Schedule планирует чтение статистики по всем провайдерам каждые every за период period.
// Возвращает буферезированный канал (длиной равной количеству запущеных провайдеров),
// откуда могут быть прочитаны очередные доступные данные статистики, готовые для отправки gRPC клиенту.
func Schedule(ctx context.Context, every time.Duration, period time.Duration) <-chan *v1.Record {
	outCh := make(chan *v1.Record, len(coordinators))
	wg := sync.WaitGroup{}

	wg.Add(len(coordinators))
	for i := range len(coordinators) {
		crd := coordinators[i]

		if agg, ok := crd.agg.(Grower); ok {
			agg.Grow(int(period / time.Second))
		}

		go func() {
			defer wg.Done()

			ch, err := crd.c.Schedule(ctx, every, period)
			if err != nil {
				crd.logger.Error("schedule failed",
					"error", err,
					"every", every,
					"period", period,
				)
				return
			}

			for v := range ch {
				if rec := crd.valueToProtoRecordFn(v); rec != nil {
					outCh <- rec
				}
			}
		}()
	}

	go func() {
		defer close(outCh)

		wg.Wait()
	}()

	return outCh
}

// initCoordinator создаёт и запускает координатор для провайдеров и агрегаторов.
// Возвращает ошибку, если запустить координатор не получилось.
//
// Основная идея:
//  1. запускаем провайдер с определённой точностью (сколько "снимков" будет сделано в секунду времени)
//  2. данные с провайдера передаются в агрегатор "provider", он расчитывает "среднее" за секунду по полученным снимкам
//  3. данные с агрегатора "provider" передаются во второй агрегатор "each second",
//     который накапливает данные каждую секунду
//
// Таким образом "выходом" каждого координатора является агрегатор "each second".
// Запланированным клиентам передаются данные с агрегатора "each second".
func initCoordinator(
	ctx context.Context,
	logger *logger.Logger,
	provider register.DataProvider,
	providerAccuracy int,
	aggMaker register.AggregatorMaker,
) (Coordinator, error) {
	crd := scheduler.NewCoordinator(logger)

	err := crd.Start(
		ctx,
		provider,
		time.Second/time.Duration(providerAccuracy),
		aggMaker(providerAccuracy),
	)
	if err != nil {
		return Coordinator{}, fmt.Errorf("can't start coordinator: %w", err)
	}

	agg := aggMaker(1)
	err = crd.AppendAggregator(
		"each second",
		agg,
		time.Second,
		time.Second,
	)
	if err != nil {
		return Coordinator{}, fmt.Errorf("can't start coordinator: %w", err)
	}

	return Coordinator{
		c:                    crd,
		logger:               logger,
		agg:                  agg,
		valueToProtoRecordFn: provider.ValueToProtoRecord,
	}, nil
}
