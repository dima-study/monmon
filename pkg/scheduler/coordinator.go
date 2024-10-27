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

type Coordinator struct {
	logger *logger.Logger

	s  *AggScheduler
	mx sync.Mutex
}

func NewCoordinator(l *logger.Logger) *Coordinator {
	coord := Coordinator{
		logger: l,
	}

	return &coord
}

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

	ch := s.s.Schedule(context.Background(), every, period)

	s.s = NewAggScheduler(
		s.logger,
		purpose,
		ch,
		agg,
	)
}

func (s *Coordinator) Schedule(
	ctx context.Context,
	every time.Duration,
	period time.Duration,
) <-chan any {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.s.Schedule(ctx, every, period)
}

func (s *Coordinator) startProvider(ctx context.Context, provider DataProvider, accuracy time.Duration) <-chan any {
	l := s.logger.With("provider", provider.String())
	l.Debug("start provider")

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
