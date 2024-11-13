package start

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dima-study/monmon/pkg/logger"
)

func unsetEnv() {
	os.Unsetenv("MONMON_SHUTDOWN_TIMEOUT")

	os.Unsetenv("MONMON_GRPC_HOST")
	os.Unsetenv("MONMON_GRPC_PORT")

	os.Unsetenv("MONMON_LOG_LEVEL")

	os.Unsetenv("MONMON_SERVICE_ACCURACY")
	os.Unsetenv("MONMON_SERVICE_MIN_INTERVAL")
	os.Unsetenv("MONMON_SERVICE_MAX_INTERVAL")
	os.Unsetenv("MONMON_SERVICE_MIN_PERIOD")
	os.Unsetenv("MONMON_SERVICE_MAX_PERIOD")
	os.Unsetenv("MONMON_SERVICE_DISABLED_PROVIDERS")
}

func Test_ParseConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  string
		init func()
		want Config
	}{
		{
			name: "full config",
			cfg: `
  shutdown_timeout: 5s

  grpc:
    port: "54321"
    host: lolo

  logger:
    level: trace

  service: 
    accuracy: 5
    min_interval: 10
    max_interval: 20
    min_period: 30
    max_period: 40
    disabled_providers:
      - loadavg
      `,
			want: Config{
				ShutdownTimeout: 5 * time.Second,

				GRPC: GRPCConfig{
					Host: "lolo",
					Port: "54321",
				},
				Log: LoggerConfig{
					Level: LogLevel(logger.LevelTrace),
				},
				Service: ServiceConfig{
					Accuracy:          5,
					MinInterval:       10,
					MaxInterval:       20,
					MinPeriod:         30,
					MaxPeriod:         40,
					DisabledProviders: []string{"loadavg"},
				},
			},
		},
		{
			name: "overwrite by env",
			cfg: `
  shutdown_timeout: 5s

  grpc:
    port: "00000"
    host: "~~~~~"

  logger:
    level: warn

  service: 
    accuracy: 1
    min_interval: 2
    max_interval: 3
    min_period: 4
    max_period: 5
    disabled_providers:
      - qrs
      - tuv
      - wxyz
      `,
			init: func() {
				os.Setenv("MONMON_SHUTDOWN_TIMEOUT", "1s")

				os.Setenv("MONMON_GRPC_HOST", "some.grpc.host")
				os.Setenv("MONMON_GRPC_PORT", "12345")

				os.Setenv("MONMON_LOG_LEVEL", "error")

				os.Setenv("MONMON_SERVICE_ACCURACY", "5")
				os.Setenv("MONMON_SERVICE_MIN_INTERVAL", "100")
				os.Setenv("MONMON_SERVICE_MAX_INTERVAL", "200")
				os.Setenv("MONMON_SERVICE_MIN_PERIOD", "300")
				os.Setenv("MONMON_SERVICE_MAX_PERIOD", "400")
				os.Setenv("MONMON_SERVICE_DISABLED_PROVIDERS", "loadavg")
			},
			want: Config{
				ShutdownTimeout: time.Second,

				GRPC: GRPCConfig{
					Host: "some.grpc.host",
					Port: "12345",
				},
				Log: LoggerConfig{
					Level: LogLevel(slog.LevelError),
				},
				Service: ServiceConfig{
					Accuracy:          5,
					MinInterval:       100,
					MaxInterval:       200,
					MinPeriod:         300,
					MaxPeriod:         400,
					DisabledProviders: []string{"loadavg"},
				},
			},
		},
		{
			name: "default",
			cfg:  `default: true`,
			want: Config{
				ShutdownTimeout: 5 * time.Second,

				GRPC: GRPCConfig{
					Host: "localhost",
					Port: "50051",
				},
				Log: LoggerConfig{
					Level: LogLevel(slog.LevelInfo),
				},
				Service: ServiceConfig{
					Accuracy:          10,
					MinInterval:       1,
					MaxInterval:       300,
					MinPeriod:         1,
					MaxPeriod:         300,
					DisabledProviders: []string{},
				},
			},
		},
	}

	for i, tt := range tests {
		name := tt.name
		if name == "" {
			name = strconv.Itoa(i)
		}

		t.Run(name, func(t *testing.T) {
			if tt.init != nil {
				tt.init()
			}

			r := strings.NewReader(tt.cfg)
			cfg, err := ParseConfig(r)

			require.NoError(t, err, "must not have error")
			require.Equal(t, tt.want, cfg, "must be equal")

			unsetEnv()
		})
	}
}

func Test_ParseConfigError(t *testing.T) {
	tests := []struct {
		name      string
		cfg       string
		init      func()
		wantError bool
	}{
		{
			name: "invalid log level",
			cfg: `
		    logger:
		      level: failed
		        `,
			wantError: true,
		},
		{
			name: "invalid log level env",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_LOG_LEVEL", "failed")
			},
			wantError: true,
		},
		{
			name: "invalid service accuracy (min)",
			cfg: `
      service:
        accuracy: -1
          `,
			wantError: true,
		},
		{
			name: "invalid service accuracy env (min)",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_ACCURACY", "0")
			},
			wantError: true,
		},
		{
			name: "invalid service min_interval (min)",
			cfg: `
      service:
        min_interval: -1
          `,
			wantError: true,
		},
		{
			name: "invalid service min_interval env (min)",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_MIN_INTERVAL", "0")
			},
			wantError: true,
		},
		{
			name: "invalid service max_interval (min)",
			cfg: `
      service:
        max_interval: -1
          `,
			wantError: true,
		},
		{
			name: "invalid service max_interval env (min)",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_MAX_INTERVAL", "0")
			},
			wantError: true,
		},
		{
			name: "invalid service min_period (min)",
			cfg: `
      service:
        min_period: -1
          `,
			wantError: true,
		},
		{
			name: "invalid service min_period env (min)",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_MIN_PERIOD", "0")
			},
			wantError: true,
		},
		{
			name: "invalid service max_period (min)",
			cfg: `
      service:
        max_period: -1
          `,
			wantError: true,
		},
		{
			name: "invalid service max_period env (min)",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_MAX_PERIOD", "0")
			},
			wantError: true,
		},
		{
			name: "invalid service min_interval<=max_interval",
			cfg: `
      service:
        min_interval: 10
        max_interval: 5
          `,
			wantError: true,
		},
		{
			name: "invalid service min_interval<=max_interval env",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_MIN_INTERVAL", "10")
				os.Setenv("MONMON_SERVICE_MAX_INTERVAL", "5")
			},
			wantError: true,
		},
		{
			name: "invalid service min_period<=max_period",
			cfg: `
      service:
        min_period: 10
        max_period: 5
          `,
			wantError: true,
		},
		{
			name: "invalid service min_period<=max_period env",
			cfg:  `default: true`,
			init: func() {
				os.Setenv("MONMON_SERVICE_MIN_PERIOD", "10")
				os.Setenv("MONMON_SERVICE_MAX_PERIOD", "5")
			},
			wantError: true,
		},
	}

	for i, tt := range tests {
		name := tt.name
		if name == "" {
			name = strconv.Itoa(i)
		}

		t.Run(name, func(t *testing.T) {
			if tt.init != nil {
				tt.init()
			}

			r := strings.NewReader(tt.cfg)
			_, err := ParseConfig(r)

			require.Error(t, err, "must have error")

			unsetEnv()
		})
	}
}
