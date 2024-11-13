package providers

import (
	"flag"
	"fmt"

	"github.com/dima-study/monmon/internal/stats/register"
)

type Cmd struct{}

func New() *Cmd {
	return &Cmd{}
}

func (cmd *Cmd) Cmd() string {
	return "providers"
}

func (cmd *Cmd) Info() string {
	return "list supported providers"
}

func (cmd *Cmd) Match(command string, args []string) bool {
	return command == cmd.Cmd()
}

func (cmd *Cmd) Run() error {
	statsList := register.SupportedStats()

	fmt.Fprintf(flag.CommandLine.Output(), "Compiled providers:\n")
	for _, s := range statsList {
		p, _ := register.GetProvider(s)
		fmt.Fprintf(flag.CommandLine.Output(), "    %s: %s\n", s, p.Name())
		if err := register.CheckStatAvailability(s); err != nil {
			fmt.Fprintf(flag.CommandLine.Output(), "        not available: %s\n", err.Error())
		}
	}

	return nil
}
