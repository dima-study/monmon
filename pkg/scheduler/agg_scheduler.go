package scheduler

import (
	"container/list"
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

	// Cleanup очищает данные агрегатора.
	// Может быть использовано для уделения временных файлов.
	Cleanup(ctx context.Context) error

	String() string
}

// AggScheduler - планировщик агрегатора данных статистики.
//
// Основная задача - из входящего потока данных от провайдера данных добавлять их в агрегатор.
// А в случае запланированного чтения (Schedule) отдавать агрегированные данные в выходной поток.
type AggScheduler struct {
	done   chan struct{}
	agg    Aggregator
	name   string
	logger *logger.Logger

	// valueAddedChs - список сигнальных каналов. Необоходимо защитить мутексом, т.к. сигнальный канал может быть закрыт
	// при завершении запланированного чтения (Schedule) и быть удалённым из списка.
	//
	// Для каждого нового запланированного чтения (Schedule) будет создан сигнльный канал.
	// Сигнальный канал может быть прочитан, если очередное значение от провайдера было добавлено в агрегатор
	// после назначения запланированного чтения.
	//
	// Назначние сигнального канала - предотвратить гонку между добавлением значения в агрегатор и чтения значения
	// из агрегатора *за период*.
	// Проблема возникает при двух "одновременных" событиях чтения и записи, когда выстраивается следующий порядок:
	//   1. провайдер данных отдал данные с меткой времени
	//   2. чтение из агрегатора данных за период
	//   3. запись в агрегатор данных с меткой времени (из п.1)
	//
	// В данном случае, когда интервал поступления данных равен интервалу чтения данных, чтение может возвращать
	// пустое значение.
	// Если интервал поступления данных меньше интервала чтения данных, чтение будет возвращать данные
	// за некорректный период.
	// Соответственно, нужно упорядочить п.2 после п.3 (т.е. сначала запись, затем чтение).
	//
	// Однако, данное решение накладывает ограничение: интервал чтения должен быть кратен интервалу записи!
	// Т.к. читатель будет ожидать очередной записи в агрегатор.
	//
	// Как вариант готовой реализации - sync.Cond, но могут быть пропуски.
	valueAddedChs *list.List
	mx            sync.RWMutex
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
		logger: logger,

		valueAddedChs: list.New(),
	}

	go func() {
		defer func() {
			s.logger.Debug("stop agg scheduler")

			close(s.done)
		}()

		for v := range ch {
			s.logger.Trace("agg.Add", "val", v)
			s.agg.Add(v)

			// Информируем "читателей", что значение добавлено в агрегатор.
			s.mx.RLock()
			for itm := s.valueAddedChs.Front(); itm != nil; itm = itm.Next() {
				select {
				case itm.Value.(chan struct{}) <- struct{}{}:
				default:
				}
			}
			s.mx.RUnlock()
		}
	}()

	return &s
}

// Wait ждёт завершения планировщика агрегатора и выполняет очистку (Cleanup) аггрегатора.
// Возвращает результат очистки.
func (s *AggScheduler) Wait(ctx context.Context) error {
	<-s.done
	s.logger.Debug("cleanup")
	return s.agg.Cleanup(ctx)
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
	s.logger.Debug("schedule",
		"every", every.String(),
		"period", period.String(),
	)

	ch := make(chan any)

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

		// Сигнальный канал
		valueAddedCh := make(chan struct{}, 1)
		defer close(valueAddedCh)

		s.mx.Lock()
		elm := s.valueAddedChs.PushFront(valueAddedCh)
		s.mx.Unlock()
		defer func() {
			s.mx.Lock()
			defer s.mx.Unlock()

			s.valueAddedChs.Remove(elm)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-s.done:
				return
			case <-t.C:
				select {
				case <-ctx.Done():
					return
				case <-s.done:
					return
				case <-valueAddedCh:
				}

				val, ok := s.agg.Get(period)
				sent := true
				if ok {
					select {
					case ch <- val:
					default:
						sent = false
					}
				} else {
					// Проблема: из-за лага получения данных провайдера, добавления в агрегатор и таймера,
					// может произойти ситуация, когда данные в агрегаторе по таймеру будут за пределами period.
					// В этой ситуации по завершению таймера в канале valueAddedCh всегда будут данные.
					//
					// Нужно "поменять местами" добавление в агрегатор и чтение из него: принудительно очищаем канал.
					// Может произойти пропуск одного значения в выдаче клиенту.
					go func() {
						<-valueAddedCh
					}()
				}

				s.logger.Trace("agg.Get",
					"ok", ok,
					"sent", sent,
					"every", every.String(),
					"period", period.String(),
					"val", val,
				)
			}
		}
	}()

	return ch
}

// String возвращает название планировщика.
func (s *AggScheduler) String() string {
	return s.name
}
