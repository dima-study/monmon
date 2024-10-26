package start

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/dima-study/monmon/pkg/logger"
)

type Config struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"MONMON_SHUTDOWN_TIMEOUT" env-default:"5s"`

	GRPC    GRPCConfig    `yaml:"grpc"    env-prefix:"MONMON_GRPC_"`
	Log     LoggerConfig  `yaml:"logger"  env-prefix:"MONMON_LOG_"`
	Service ServiceConfig `yaml:"service" env-prefix:"MONMON_SERVICE_"`
}

func (c *Config) Validate() error {
	if err := c.Service.Validate(); err != nil {
		return err
	}

	return nil
}

type GRPCConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"PORT" env-default:"50051"`
}

type LoggerConfig struct {
	Level LogLevel `yaml:"level" env:"LEVEL" env-default:"info"`
}

type LogLevel slog.Level

func (l *LogLevel) UnmarshalText(s []byte) error {
	var ll slog.Level
	if err := ll.UnmarshalText(s); err == nil {
		*l = LogLevel(ll)
		return nil
	}

	switch strings.ToUpper(string(s)) {
	case "TRACE":
		*l = LogLevel(logger.LevelTrace)
	default:
		return fmt.Errorf("invalid log level '%s'", s)
	}

	return nil
}

func (l *LogLevel) String() string {
	if slog.Level(*l) == logger.LevelTrace {
		return "TRACE"
	}

	return slog.Level(*l).String()
}

type ServiceConfig struct {
	// Accuracy - сколько делать в секунду "снимков" статистики по каждому провайдеру.
	Accuracy int `yaml:"accuracy" env:"ACCURACY" env-default:"10"`

	// Min(Max)Interval мин/макс интервал для получения статистики по запросам.
	MinInterval int `yaml:"min_interval" env:"MIN_INTERVAL" env-default:"1"`
	MaxInterval int `yaml:"max_interval" env:"MAX_INTERVAL" env-default:"300"`

	// Min(Max)Period мин/макс период за который возможно получить статистику.
	MinPeriod int `yaml:"min_period" env:"MIN_PERIOD" env-default:"1"`
	MaxPeriod int `yaml:"max_period" env:"MAX_PERIOD" env-default:"300"`

	// Список отключённых провайдеров статистики.
	DisabledProviders []string `yaml:"disabled_providers" env:"DISABLED_PROVIDERS" env-default:""`
}

func (sc *ServiceConfig) Validate() error {
	posIntChecks := []struct {
		n string
		v int
	}{
		{"accuracy", sc.Accuracy},
		{"min_interval", sc.MinInterval},
		{"max_interval", sc.MaxInterval},
		{"min_period", sc.MinPeriod},
		{"max_period", sc.MaxPeriod},
	}

	for _, v := range posIntChecks {
		if v.v <= 0 {
			return fmt.Errorf("invalid %s (%d): %s must be >= 0", v.n, v.v, v.n)
		}
	}

	if sc.MinInterval > sc.MaxInterval {
		return fmt.Errorf(
			"min_interval > max_interval (%d > %d): min_interval must not be great than max_interval",
			sc.MinInterval, sc.MaxInterval,
		)
	}

	if sc.MinPeriod > sc.MaxPeriod {
		return fmt.Errorf(
			"min_period > max_period (%d > %d): min_period must not be great than max_period",
			sc.MinPeriod, sc.MaxPeriod,
		)
	}

	return nil
}

// ReadConfig пытается прочитать конфиг в yaml формате из файла и переменных окружения.
func ReadConfig(path string) (Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// ParseConfig пытается прочитать конфиг в yaml формате из r и переменных окружения.
func ParseConfig(r io.Reader) (Config, error) {
	var cfg Config
	if err := cleanenv.ParseYAML(r, &cfg); err != nil {
		return Config{}, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
