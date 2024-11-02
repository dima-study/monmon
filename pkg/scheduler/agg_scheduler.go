package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/dima-study/monmon/pkg/logger"
)

// Aggregator представляет собой агрегатор статистики.
type Aggregator interface {
	// Add добавляет данные в агрегатор
	Add(val any)

	// Get получает усреднённые агрегированные данные за период period
	Get(period time.Duration) (any, bool)

	String() string
}

// AggScheduler - планировщик агрегатора данных статистики.
//
// Основная задача - из входящего потока данных от провайдера данных добавлять их в агрегатор.
// А в случае запланированного чтения отдавать агрегированные данные в выходной поток.
type AggScheduler struct {
	done   chan struct{}
	agg    Aggregator
	name   string
	cond   *sync.Cond
	logger *logger.Logger
}

// NewAggScheduler создаёт новый планировщик агрегатора:
// добавляет данные из входящего потока ch от провайдера данных в агрегатор agg.
//
// Закрытие входящий потока ch завершает работу планировщика.
func NewAggScheduler(logger *logger.Logger, name string, ch <-chan any, agg Aggregator) *AggScheduler {
	logger = logger.With(
		"agg_scheduler", name,
		"aggregator", agg.String(),
	)

	s := AggScheduler{
		done:   make(chan struct{}),
		agg:    agg,
		name:   name,
		cond:   sync.NewCond(&sync.Mutex{}),
		logger: logger,
	}

	go func() {
		defer func() {
			s.logger.Debug("stop agg scheduler")

			close(s.done)

			s.cond.L.Lock()
			s.cond.Broadcast()
			s.cond.L.Unlock()
		}()

		for v := range ch {
			s.cond.L.Lock()

			s.logger.Trace("agg.Add", "val", v)
			s.agg.Add(v)

			s.cond.Broadcast()
			s.cond.L.Unlock()
		}
	}()

	return &s
}

// Wait ждёт завершения планировщика агрегатора.
func (s *AggScheduler) Wait() {
	<-s.done
}

// Schedule создаёт канал запланированного чтения: через каждый every интервал в канал будут переданы
// усреднённые агрегированные данные за период period.
//
// При завершении работы планировщика передача данных в каналы запланированных чтениий будет также завершена,
// а сами каналы будут закрыты.
//
// Завершить определённое запланированное чтение возможно через отмену контекста ctx.
//
// В случае, когда данные не могут быть записаны в канал - они будут утеряны.
func (s *AggScheduler) Schedule(ctx context.Context, every time.Duration, period time.Duration) <-chan any {
	ch := make(chan any)

	s.logger.Debug("schedule",
		"every", every.String(),
		"period", period.String(),
	)

	go func() {
		defer func() {
			s.logger.Debug("stop agg scheduling",
				"every", every.String(),
				"period", period.String(),
			)

			close(ch)
		}()

		t := time.NewTicker(every)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-s.done:
				return
			default:
			}

			s.cond.L.Lock()
			s.cond.Wait()

			select {
			case <-t.C:
				val, ok := s.agg.Get(period)
				sent := true
				if ok {
					select {
					case ch <- val:
					default:
						sent = false
					}
				}

				s.logger.Trace("agg.Get",
					"ok", ok,
					"sent", sent,
					"every", every.String(),
					"period", period.String(),
					"val", val,
				)
			default:
			}

			s.cond.L.Unlock()
		}
	}()

	return ch
}

// String возвращает название планировщика.
func (s *AggScheduler) String() string {
	return s.name
}
