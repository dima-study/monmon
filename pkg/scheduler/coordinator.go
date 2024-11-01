package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dima-study/monmon/pkg/logger"
)

var ErrAlreadyStarted = errors.New("coordinator is already started")

type DataProvider interface {
	// Available возвращает ошибку по которой провайдер не доступен
	Available() error

	// Data возвращает данные, которые предоставляет провайдер.
	Data() (any, error)

	String() string
}

// Coordinator коодринатор планировщика.
// Задача - связать провайдера статистики, агрегаторы и планировать входящие запросы на получение статистики.
type Coordinator struct {
	logger *logger.Logger

	s  *AggScheduler
	mx sync.Mutex
}

// NewCoordinator создаёт новый координатор.
func NewCoordinator(l *logger.Logger) *Coordinator {
	coord := Coordinator{
		logger: l,
	}

	return &coord
}

// Start запускает кординатор для указанного агрегатора agg и провайдера provider с указанной точностью accuracy.
// Если координатор c был запущен ранее, будет возвращена ошибка ErrAlreadyStarted.
// Если провайдер не доступен, будет возвращена ошибка с причиной.
func (s *Coordinator) Start(
	ctx context.Context,
	provider DataProvider,
	accuracy time.Duration,
	agg Aggregator,
) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.s != nil {
		return ErrAlreadyStarted
	}

	if err := provider.Available(); err != nil {
		return fmt.Errorf("%s: %w", provider.String(), err)
	}

	s.s = NewAggScheduler(
		s.logger,
		"stat provider",
		s.startProvider(ctx, provider, accuracy),
		agg,
	)

	return nil
}

// AppendAggregator позволяет собирать цепочку из агрегаторов.
// Решает ситуацию, когда первый/входящий агрегатор собирает более точные данные/более часто, чем последний/выходной.
func (s *Coordinator) AppendAggregator(
	purpose string,
	agg Aggregator,
	every time.Duration,
	period time.Duration,
) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.s == nil {
		panic("coordinator is not started")
	}

	s.logger.Debug("append aggregator",
		"aggregator", agg.String(),
		"append", purpose,
		"to", s.s.String(),
		"every", every.String(),
		"period", period.String(),
	)

	// context.Background т.к запланированное чтение будет завершено при завершении родительского планировщика.
	ch := s.s.Schedule(context.Background(), every, period)

	s.s = NewAggScheduler(
		s.logger,
		purpose,
		ch,
		agg,
	)
}

// Schedule возвращает канал (буферезированный на 1 элемент) с данными запланированного чтения из выходного агрегатора
// каждые every за период period.
func (s *Coordinator) Schedule(
	ctx context.Context,
	every time.Duration,
	period time.Duration,
) <-chan any {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.s.Schedule(ctx, every, period)
}

// startProvider запускает провайдер статистики и возвращает канал с данными от провайдера.
// Данные от провайдера направляются в буферезированный канал (1 элемент).
// Если данные не могут быть записаны в канал, будет логировано Warn-сообщение о неудаче.
// Если данные не могут быть получены от провайдера, будет логировано Error-сообщение о неудаче.
func (s *Coordinator) startProvider(ctx context.Context, provider DataProvider, accuracy time.Duration) <-chan any {
	l := s.logger.With("provider", provider.String())
	l.Debug("start provider", "accuracy", accuracy)

	ch := make(chan any, 1)

	go func() {
		defer l.Debug("stop provider")
		defer close(ch)

		tic := time.NewTicker(accuracy)
		defer tic.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case <-tic.C:
				v, err := provider.Data()
				if err != nil {
					l.Error("provider.Data", "error", err)
					continue
				}

				select {
				case ch <- v:
				default:
					l.Warn("value is not sent")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}
