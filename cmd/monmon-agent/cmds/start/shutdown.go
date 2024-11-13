package start

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dima-study/monmon/pkg/logger"
)

type (
	StartServerFunc func() error //nolint:revive
	StopServerFunc  func(context.Context) error
)

func StartAndShutdown( //nolint:revive
	ctx context.Context,
	logger *logger.Logger,
	start StartServerFunc,
	stop StopServerFunc,
	stopTimeout time.Duration,
) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	startErrCh := make(chan error)
	go func() {
		if err := start(); err != nil {
			startErrCh <- err
		}
	}()

	select {
	case err := <-startErrCh:
		return fmt.Errorf("server start error: %w", err)
	case sig := <-shutdown:
		logger.Info(
			"shutdown",
			slog.String("signal", sig.String()),
		)

		ctx, cancel := context.WithTimeout(ctx, stopTimeout)
		defer cancel()

		if err := stop(ctx); err != nil {
			return fmt.Errorf("server stop error: %w", err)
		}
	}

	return nil
}
