package main

import (
	"fmt"
	"strings"

	core "github.com/HeminWon/proteus/internal/cli"
)

var switchFlags = []string{"--dry-run", "--help", "-h"}

func parseSwitchArgs(args []string) (core.CliOptions, error) {
	dryRun := false
	help := false
	positional := make([]string, 0)

	for _, arg := range args {
		switch {
		case arg == "--dry-run":
			dryRun = true
		case isHelpFlag(arg):
			help = true
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, unknownOptionError(string(core.ActionSwitch), arg, switchFlags)
		default:
			positional = append(positional, arg)
		}
	}

	if help {
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with switch --help")
		}
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("switch --help cannot be combined with other options")
		}
		return core.CliOptions{Action: core.ActionHelp, HelpCommand: string(core.ActionSwitch)}, nil
	}

	if len(positional) == 0 {
		return core.CliOptions{}, fmt.Errorf("missing provider for switch")
	}
	if len(positional) > 1 {
		return core.CliOptions{}, fmt.Errorf("too many provider arguments: %s", strings.Join(positional, ", "))
	}

	return core.CliOptions{Action: core.ActionSwitch, ProviderInput: positional[0], DryRun: dryRun}, nil
}
