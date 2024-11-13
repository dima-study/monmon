package scheduler

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/dima-study/monmon/pkg/logger"
)

func TestCoordinator(t *testing.T) {
	suite.Run(t, new(CoordinatorTestSuite))
}

type FakeProvider struct {
	available bool
	n         int
	cleaned   bool
	mx        sync.Mutex
}

var errNotAvailable = errors.New("not available")

func (p *FakeProvider) Available() error {
	if !p.available {
		return errNotAvailable
	}
	return nil
}

func (p *FakeProvider) Data() (any, error) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.n++
	return p.n, nil
}

func (p *FakeProvider) N() int {
	p.mx.Lock()
	defer p.mx.Unlock()

	return p.n
}

func (p *FakeProvider) Cleanup(ctx context.Context) error {
	p.cleaned = true
	return nil
}

func (p *FakeProvider) String() string {
	return "FakeProvider"
}

type CoordinatorTestSuite struct {
	suite.Suite

	c   *Coordinator
	p   *FakeProvider
	agg *FakeAgg
}

func (s *CoordinatorTestSuite) SetupTest() {
	h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: logger.LevelTrace})
	l := logger.New(h)

	s.c = NewCoordinator(l)
	s.p = &FakeProvider{available: true}
	s.agg = &FakeAgg{}
}

func (s *CoordinatorTestSuite) TestStart() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.p.available = false
	err := s.c.Start(ctx, s.p, 100*time.Millisecond, s.agg)
	s.Require().ErrorIs(err, errNotAvailable, "must be errNotAvailable error")

	s.p.available = true
	err = s.c.Start(ctx, s.p, 100*time.Millisecond, s.agg)
	s.Require().NoError(err, "must be no error")
	defer s.c.Reset(context.Background())

	err = s.c.Start(ctx, s.p, 100*time.Millisecond, s.agg)
	s.Require().ErrorIs(err, ErrAlreadyStarted, "must be ErrAlreadyStarted error")
}

func (s *CoordinatorTestSuite) TestAppendAggregator() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.c.AppendAggregator("second", &FakeAgg{}, 100*time.Millisecond, 100*time.Millisecond)
	s.Require().ErrorIs(err, ErrNotStarted, "must be ErrNotStarted error")

	err = s.c.Start(ctx, s.p, 100*time.Millisecond, s.agg)
	s.Require().NoError(err, "must be no error")
	defer s.c.Reset(context.Background())

	err = s.c.AppendAggregator("second", &FakeAgg{}, 100*time.Millisecond, 100*time.Millisecond)
	s.Require().NoError(err, "must be no error")
}

func (s *CoordinatorTestSuite) TestAppendSchedule() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := s.c.Start(ctx, s.p, 100*time.Millisecond, s.agg)
	s.Require().NoError(err, "must be no error")
	defer s.c.Reset(context.Background())

	agg11 := FakeAgg{m: 11}
	err = s.c.AppendAggregator("second", &agg11, 100*time.Millisecond, 100*time.Millisecond)
	s.Require().NoError(err, "must be no error")

	sCtx, sCancel := context.WithCancel(ctx)
	defer sCancel()
	ch, err := s.c.Schedule(sCtx, 100*time.Millisecond, 100*time.Millisecond)
	s.Require().NoError(err, "must be no error")

	got := []int{}
	for v := range ch {
		i := v.(int)
		got = append(got, i)
		if len(got) == 5 {
			sCancel()
		}
	}

	<-ctx.Done()

	s.Require().True(5 < s.p.N() && s.p.N() <= 10, "must sent >5 elements")
	s.Require().Len(got, 5, "result must have 5 elements")
	for _, i := range got {
		s.Require().True(i%11 == 0, "must be devidable")
	}
}

func (s *CoordinatorTestSuite) TestAppendReset() {
	try := func(name string) {
		s.Run(name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			err := s.c.Start(ctx, s.p, 100*time.Millisecond, s.agg)
			s.Require().NoError(err, "must be no error")
			defer s.c.Reset(context.Background())

			agg11 := FakeAgg{m: 11}
			err = s.c.AppendAggregator("second", &agg11, 100*time.Millisecond, 100*time.Millisecond)
			s.Require().NoError(err, "must be no error")

			sCtx, sCancel := context.WithCancel(ctx)
			defer sCancel()
			ch, err := s.c.Schedule(sCtx, 100*time.Millisecond, 100*time.Millisecond)
			s.Require().NoError(err, "must be no error")

			got := []int{}
			for v := range ch {
				i := v.(int)
				got = append(got, i)
				if len(got) == 5 {
					s.c.Reset(context.Background())
				}
			}

			<-ctx.Done()

			s.Require().Len(got, 5, "must have 5 elements")
			for _, i := range got {
				s.Require().True(i%11 == 0, "must be devidable")
			}

			s.Require().True(s.p.cleaned, "provider must be cleaned")
			s.Require().True(s.agg.cleaned, "first aggregator must be cleaned")
			s.Require().True(agg11.cleaned, "second agg must be cleaned")
		})
	}

	try("first start and reset")
	try("second start")
}
