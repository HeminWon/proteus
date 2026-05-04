package main

import (
	"fmt"
	"strings"

	core "github.com/HeminWon/proteus/internal/cli"
)

func parseLaunchArgs(args []string) (core.CliOptions, error) {
	dryRun := false
	list := false
	help := false
	positional := make([]string, 0)

	for _, arg := range args {
		switch {
		case arg == "--dry-run":
			dryRun = true
		case arg == "--list":
			list = true
		case isHelpFlag(arg):
			help = true
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, fmt.Errorf("unknown launch option: %s%s", arg, core.SuggestFlag(arg, []string{"--dry-run", "--list", "--help", "-h"}))
		default:
			positional = append(positional, arg)
		}
	}

	if help {
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected profile argument with launch --help")
		}
		if list || dryRun {
			return core.CliOptions{}, fmt.Errorf("launch --help cannot be combined with other options")
		}
		return core.CliOptions{Action: core.ActionHelp, HelpCommand: string(core.ActionLaunch)}, nil
	}

	if list {
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected profile argument with launch --list")
		}
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run cannot be used with launch --list")
		}
		return core.CliOptions{Action: core.ActionLaunch, ListLaunch: true}, nil
	}

	if len(positional) == 0 {
		return core.CliOptions{}, fmt.Errorf("missing profile for launch")
	}
	if len(positional) > 1 {
		return core.CliOptions{}, fmt.Errorf("too many profile arguments: %s", strings.Join(positional, ", "))
	}

	return core.CliOptions{Action: core.ActionLaunch, ProfileInput: positional[0], DryRun: dryRun}, nil
}
