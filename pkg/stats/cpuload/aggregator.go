package cpuload

import (
	"context"
	"sync"
	"time"
)

// CPULoad - агрегатор для расчёта среднего средней загрузки ЦПУ.
type CPULoad struct {
	idx  int     // индекс *следующего* элемента, который будет добавлен в буфер
	buf  []Value // кольцевой буфер, хранит среднее значение на указанное время Value.T
	prec int     // точность расчёта среднего

	mx sync.RWMutex
}

// NewAggregator создаёт новый агрегатор CPULoad на n значений.
// n должно быть больше 0, иначе возникнет паника!
func NewAggregator(n int) *CPULoad {
	if n <= 0 {
		panic("len must be >=0")
	}

	return &CPULoad{
		// +1 элемент для "кеша":
		// когда всего один элемент будет частая гонка чтения-записи за него
		buf:  make([]Value, n+1),
		prec: precision(n),
	}
}

// Add добавляет значение в агрегатор и высчитывает среднее значение по каждому срезу времени.
func (agg *CPULoad) Add(val Value) {
	// Среднее будет рассчитано с учётом точности.
	val.User = val.User * agg.prec
	val.System = val.System * agg.prec
	val.Idle = val.Idle * agg.prec

	// Добавление и расчёт среднего значения происходит в локе
	agg.mx.Lock()
	defer agg.mx.Unlock()

	// Не выходим за границы кольцевого буфера.
	i := agg.idx
	agg.idx++
	if agg.idx == len(agg.buf) {
		agg.idx = 0
	}

	agg.buf[i] = val

	n := 1

	// Для каждого элемента в буфере пересчитываем среднее на основании нового добавленного значения
	for range len(agg.buf) {
		i--
		if i < 0 {
			i = len(agg.buf) - 1
		}

		// Достигли начала буфера (буфер заполнен не полностью)
		if agg.buf[i].IsEmpty() {
			break
		}

		// Вышли либо на себя либо на значение в прошлом с некорректным временем в будущем (относительно нового значения).
		//
		// НЕ ДОЛЖНО случиться!
		// Может случиться, если val передан с некорректным времене T.
		if !agg.buf[i].T.Before(val.T) {
			break
		}

		// Считааем среднее
		agg.buf[i].User = (agg.buf[i].User*n + val.User) / (n + 1)
		agg.buf[i].System = (agg.buf[i].System*n + val.System) / (n + 1)
		agg.buf[i].Idle = (agg.buf[i].Idle*n + val.Idle) / (n + 1)

		n++
	}
}

// Get возвращает среднее значение за указанный период (period) в прошлом относительно текущего времени.
// Время Т значения будет соответствовать времени на которое среднее значение расчитано.
func (agg *CPULoad) Get(period time.Duration) (Value, bool) {
	agg.mx.RLock()
	defer agg.mx.RUnlock()

	buf := agg.buf
	prec := agg.prec
	i := agg.idx - 1

	if i < 0 {
		i = len(buf) - 1
	}

	t := time.Now().Add(-period)

	avg := Value{}
	avgT := buf[i].T

	// Идём по каждому элементу в буфере пока не выйдем за период
	// -1 потому что буфер содержит "лишний" элемент (см. New)
	for range len(buf) - 1 {
		if i < 0 {
			i = len(buf) - 1
		}

		val := buf[i]

		// Буфер заполнен не полностью, не нашли данных за период
		if val.IsEmpty() {
			return Value{}, false
		}

		// Вышли за период
		if val.T.Before(t) {
			break
		}

		avg = val

		i--
	}

	// Не нашли данных для заданного периода
	if avg.IsEmpty() {
		return Value{}, false
	}

	// Возвращаем значение (без точности)
	avg.User = avg.User / prec
	avg.System = avg.System / prec
	avg.Idle = avg.Idle / prec
	avg.T = avgT

	return avg, true
}

// Grow увеличивает размер хранилища до n элементов по которым считать среднее значение.
func (agg *CPULoad) Grow(n int) {
	agg.mx.Lock()
	defer agg.mx.Unlock()

	if n <= len(agg.buf)-1 {
		return
	}

	n = n - len(agg.buf) + 1

	prec := precision(len(agg.buf) + n)
	if prec != agg.prec {
		mul := prec / agg.prec

		for i := range len(agg.buf) {
			agg.buf[i].User = mul * agg.buf[i].User
			agg.buf[i].System = mul * agg.buf[i].System
			agg.buf[i].Idle = mul * agg.buf[i].Idle
		}

		agg.prec = prec
	}

	apnd := make([]Value, n)
	agg.buf = append(agg.buf[:agg.idx], append(apnd, agg.buf[agg.idx:]...)...) //nolint:makezero
}

func (agg *CPULoad) String() string {
	return "cpuload"
}

func (agg *CPULoad) Cleanup(ctx context.Context) error {
	return nil
}

// precision возвращает необходимую точность для количества элементов в буфере:
//   - 1 для 1-9
//   - 10 для 10-99
//   - 100 для 100-999
//     ...
func precision(n int) int {
	p := 1
	for n > 0 {
		p = p * 10
		n = n / 10
	}

	return p
}
