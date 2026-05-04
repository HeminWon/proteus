package main

import (
	"fmt"
	"strings"

	core "github.com/HeminWon/proteus/internal/cli"
)

const (
	commandList     = "list"
	globalHelpAlias = "global"
)

var canonicalCommands = []string{commandList, string(core.ActionValidate), string(core.ActionSwitch), string(core.ActionLaunch)}

func isHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h"
}

func isKnownCommand(cmd string) bool {
	for _, c := range canonicalCommands {
		if cmd == c {
			return true
		}
	}
	return false
}

func unknownOptionError(command, option string, candidates []string) error {
	return fmt.Errorf("unknown %s option: %s%s", command, option, core.SuggestFlag(option, candidates))
}

func parseArgs(argv []string) (core.CliOptions, error) {
	if len(argv) == 0 {
		return core.CliOptions{Action: core.ActionHelp}, nil
	}

	cmd := argv[0]
	rest := argv[1:]
	options := []string{"--help", "-h"}

	switch cmd {
	case commandList:
		if len(rest) == 0 {
			return core.CliOptions{Action: core.ActionList}, nil
		}
		if len(rest) == 1 && isHelpFlag(rest[0]) {
			return core.CliOptions{Action: core.ActionHelp, HelpCommand: commandList}, nil
		}
		return core.CliOptions{}, fmt.Errorf("list does not accept arguments")
	case string(core.ActionValidate):
		return parseValidateArgs(rest)
	case string(core.ActionSwitch):
		return parseSwitchArgs(rest)
	case string(core.ActionLaunch):
		return parseLaunchArgs(rest)
	case "--help", "-h":
		if len(rest) == 0 {
			return core.CliOptions{Action: core.ActionHelp}, nil
		}
		if len(rest) == 1 {
			sub := rest[0]
			if isKnownCommand(sub) {
				return core.CliOptions{Action: core.ActionHelp, HelpCommand: sub}, nil
			}
			return core.CliOptions{}, fmt.Errorf("unknown command: %s%s", sub, core.SuggestCommand(sub, canonicalCommands))
		}
		return core.CliOptions{}, fmt.Errorf("too many arguments with --help")
	default:
		if strings.HasPrefix(cmd, "-") {
			return core.CliOptions{}, fmt.Errorf("unknown option: %s%s", cmd, core.SuggestFlag(cmd, options))
		}
		return core.CliOptions{}, fmt.Errorf("unknown command: %s%s", cmd, core.SuggestCommand(cmd, canonicalCommands))
	}
}
