package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	cmd := parseCmd()

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Run command %s failed: %s\n", cmd.Cmd(), err)
		os.Exit(1)
	}
}

type cmdRunner interface {
	Cmd() string
	Info() string
	Match(string, []string) bool
	Run() error
}

func parseCmd() cmdRunner {
	cmds := []cmdRunner{}

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

	if cmd == nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s COMMAND\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Available commands:\n")
		for _, r := range cmds {
			fmt.Fprintf(flag.CommandLine.Output(), "    %s\n", r.Cmd())
			fmt.Fprintf(flag.CommandLine.Output(), "        %s\n", r.Info())
		}

		os.Exit(1)
	}

	return cmd
}
