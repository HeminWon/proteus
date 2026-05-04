package main

import (
	"fmt"
	"strings"

	core "github.com/HeminWon/proteus/internal/cli"
)

func isHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h"
}

func parseArgs(argv []string) (core.CliOptions, error) {
	if len(argv) == 0 {
		return core.CliOptions{Action: core.ActionHelp}, nil
	}

	cmd := argv[0]
	rest := argv[1:]
	options := []string{"--help", "-h"}

	switch cmd {
	case "list":
		if len(rest) == 0 {
			return core.CliOptions{Action: core.ActionList}, nil
		}
		if len(rest) == 1 && isHelpFlag(rest[0]) {
			return core.CliOptions{Action: core.ActionHelp, HelpCommand: "list"}, nil
		}
		return core.CliOptions{}, fmt.Errorf("list does not accept arguments")
	case "validate":
		return parseValidateArgs(rest)
	case "switch":
		return parseSwitchArgs(rest)
	case "launch":
		return parseLaunchArgs(rest)
	case "--help", "-h":
		if len(rest) == 0 {
			return core.CliOptions{Action: core.ActionHelp}, nil
		}
		if len(rest) == 1 {
			sub := rest[0]
			if sub == "list" || sub == string(core.ActionSwitch) || sub == string(core.ActionLaunch) || sub == string(core.ActionValidate) {
				return core.CliOptions{Action: core.ActionHelp, HelpCommand: sub}, nil
			}
			return core.CliOptions{}, fmt.Errorf("unknown command: %s%s", sub, core.SuggestCommand(sub, []string{"list", "validate", "switch", "launch"}))
		}
		return core.CliOptions{}, fmt.Errorf("too many arguments with --help")
	default:
		if strings.HasPrefix(cmd, "-") {
			return core.CliOptions{}, fmt.Errorf("unknown option: %s%s", cmd, core.SuggestFlag(cmd, options))
		}
		return core.CliOptions{}, fmt.Errorf("unknown command: %s%s", cmd, core.SuggestCommand(cmd, []string{"list", "validate", "switch", "launch"}))
	}
}
