package start

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/dima-study/monmon/cmd/monmon-agent/build"
	"github.com/dima-study/monmon/internal/stats/register"
	"github.com/dima-study/monmon/pkg/logger"
)

type Cmd struct {
	configFile string
}

func New() *Cmd {
	return &Cmd{}
}

func (cmd *Cmd) Cmd() string {
	return "start"
}

func (cmd *Cmd) Info() string {
	return "start monmon agent"
}

func (cmd *Cmd) Match(command string, args []string) bool {
	return command == cmd.Cmd()
}

func (cmd *Cmd) Run() error {
	os.Args = append(os.Args[0:1], os.Args[2:]...)
	cmd.initFlag()

	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo)

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       levelVar,
		ReplaceAttr: logger.ReplaceAttrLevel,
	})

	logger := logger.New(h)

	return run(context.Background(), logger, levelVar, cmd.configFile)
}

func (cmd *Cmd) initFlag() {
	flag.StringVar(&cmd.configFile, "config", "monmon.yaml", "Path to configuration file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s %s:\n", os.Args[0], cmd.Cmd())
		flag.PrintDefaults()

		var cfg Config
		help, _ := cleanenv.GetDescription(&cfg, nil)
		fmt.Fprintf(flag.CommandLine.Output(), "\n%s\n", help)
	}

	flag.Parse()
}

func run(ctx context.Context, logger *logger.Logger, levelVar *slog.LevelVar, configFile string) (err error) {
	logger.Info(
		"starting agent",
		slog.Group(
			"build",
			"release", build.Release,
			"date", build.Date,
			"gitHash", build.GitHash,
		),
	)

	logger.Info(
		"read config",
		"file", configFile,
	)
	cfg, err := ReadConfig(configFile)
	if err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	logger.Info(
		"set logger level",
		"from", levelVar.Level().String(),
		"to", cfg.Log.Level.String(),
	)
	levelVar.Set(slog.Level(cfg.Log.Level))

	for _, p := range cfg.Service.DisabledProviders {
		logger.Info(
			"disable provider",
			"providerID", p,
		)

		if err := register.DisableStat(p); err != nil {
			return fmt.Errorf("can't disable provider %s: %w", p, err)
		}
	}

	// Контекст для координаторов
	coordCtx, cancelCoord := context.WithCancel(ctx)
	defer func() {
		// Если случились ошибки при завершении работы координаторов, добавляем их в возврат.
		errors.Join(err, ResetCoordinators(ctx))
	}()
	defer cancelCoord()

	if err := InitCoordinators(coordCtx, logger, cfg.Service.Accuracy); err != nil {
		return fmt.Errorf("can't init coordinators: %w", err)
	}

	startServer, stopServer := CreateServer(logger, &cfg)

	// Отменяем контекст координаторов, после чего гасим сервер.
	stopService := func(ctx context.Context) error {
		cancelCoord()

		return stopServer(ctx)
	}

	return StartAndShutdown(ctx, logger, startServer, stopService, cfg.ShutdownTimeout)
}
