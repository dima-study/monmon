package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dima-study/monmon/cmd/monmon-agent/cmds/providers"
	"github.com/dima-study/monmon/cmd/monmon-agent/cmds/start"
)

func main() {
	cmd := parseCmd()
	if cmd == nil {
		os.Exit(1)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Run command %s failed: %s\n", cmd.Cmd(), err)
		os.Exit(1)
	}
}

// cmdRunner интефрейс команд приложения
type cmdRunner interface {
	// Cmd возвращает название команды
	Cmd() string

	// Info возвращает информацию о команде
	Info() string

	// Match возвращает true, если данная команда является cmd и поддерживает переданные параметры
	Match(cmd string, args []string) bool

	// Run запускает команду, возвращает ошибку выполнения команды
	Run() error
}

// parseCmd возвращает необходимую команду для выполнения.
func parseCmd() cmdRunner {
	cmds := []cmdRunner{providers.New(), start.New()}

	// Находим необходимую команду.
	var cmd cmdRunner
	if len(os.Args) > 1 {
		for _, c := range cmds {
			if c.Match(os.Args[1], os.Args[2:]) {
				cmd = c
				break
			}
		}

		if cmd == nil {
			fmt.Fprintf(flag.CommandLine.Output(), "Invalid command: %s\n\n", os.Args[1])
		}
	}

	// Если команда не найдена, печатаем сообщение о том как правильно запускать приложение.
	if cmd == nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s COMMAND\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Available commands:\n")
		for _, r := range cmds {
			fmt.Fprintf(flag.CommandLine.Output(), "    %s\n", r.Cmd())
			fmt.Fprintf(flag.CommandLine.Output(), "        %s\n", r.Info())
		}
	}

	return cmd
}
