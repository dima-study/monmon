package scheduler

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dima-study/monmon/pkg/logger"
)

type FakeAgg struct {
	m       int
	values  []int
	mx      sync.Mutex
	cleaned bool
}

func NewFakeAgg() *FakeAgg {
	a := &FakeAgg{}
	return a
}

func (a *FakeAgg) Add(val any) {
	switch v := val.(type) {
	case int:
		a.mx.Lock()
		if a.m > 0 {
			v = v * a.m
		}
		a.values = append(a.values, v)
		a.mx.Unlock()
	default:
	}
}

func (a *FakeAgg) Get(period time.Duration) (any, bool) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if len(a.values) == 0 {
		return nil, false
	}

	return a.values[len(a.values)-1], true
}

func (a *FakeAgg) String() string {
	return "FakeAgg"
}

func (a *FakeAgg) Cleanup(ctx context.Context) error {
	a.cleaned = true
	return nil
}

func newAggScheduler(t *testing.T, name string) (chan any, *FakeAgg, *AggScheduler) {
	t.Helper()

	h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: logger.LevelTrace})
	l := logger.New(h)

	ch := make(chan any)
	agg := NewFakeAgg()
	sch := NewAggScheduler(l, name, ch, agg)

	return ch, agg, sch
}

func TestAggScheduler(t *testing.T) {
	providerCh, agg, sch := newAggScheduler(t, "test")

	for i := range 10 {
		providerCh <- i
	}

	close(providerCh)
	sch.Wait(context.Background())

	require.Len(t, agg.values, 10, "must be correct size")
	for i := range 10 {
		require.Equalf(t, agg.values[i], i, "[%d] %d must equal %d", i, i, agg.values[i])
	}
}

func TestAggScheduler_Schedule(t *testing.T) {
	providerCh, agg, sch := newAggScheduler(t, "test")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	var got []int
	go func() {
		defer wg.Done()

		aggCh := sch.Schedule(ctx, 100*time.Millisecond, 100*time.Millisecond)
		for v := range aggCh {
			got = append(got, v.(int))
		}
	}()

	var send []int
	go func() {
		defer close(providerCh)

		tkr := time.NewTicker(100 * time.Millisecond)
		defer tkr.Stop()

		n := 0
		for {
			n++

			select {
			case <-ctx.Done():
				return
			case <-tkr.C:
				providerCh <- n
				send = append(send, n)
			}
		}
	}()

	wg.Wait()
	sch.Wait(context.Background())

	// t.Log(send)
	// t.Log(got)

	require.Equal(t, send, agg.values, "must agg all sent")
	require.Equal(t, send[:len(got)], got, "must receive all but last")
}

func TestAggScheduler_String(t *testing.T) {
	const name = "scheduler name"
	providerCh, _, sch := newAggScheduler(t, name)

	close(providerCh)
	sch.Wait(context.Background())

	require.Equal(t, name, sch.String(), "must return correct name")
}
