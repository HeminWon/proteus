package cli

import (
	"fmt"
	"strings"

	"github.com/HeminWon/proteus/go/internal/core"
)

func ParseArgs(argv []string) (core.CliOptions, error) {
	dryRun := false
	action := ""
	positional := make([]string, 0)

	for _, arg := range argv {
		switch {
		case arg == "--dry-run":
			dryRun = true
		case arg == "--help" || arg == "-h":
			if action != "" && action != string(core.ActionHelp) {
				return core.CliOptions{}, fmt.Errorf("--help cannot be combined with other actions")
			}
			action = string(core.ActionHelp)
		case arg == "--list":
			if action != "" && action != string(core.ActionList) {
				return core.CliOptions{}, fmt.Errorf("--list cannot be combined with other actions")
			}
			action = string(core.ActionList)
		case arg == "--validate" || arg == "validate":
			if action != "" && action != string(core.ActionValidate) {
				return core.CliOptions{}, fmt.Errorf("--validate cannot be combined with other actions")
			}
			action = string(core.ActionValidate)
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, fmt.Errorf("Unknown option: %s", arg)
		default:
			positional = append(positional, arg)
		}
	}

	if action == string(core.ActionHelp) {
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run cannot be used with --help")
		}
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("Unexpected provider argument with --help")
		}
		return core.CliOptions{Action: core.ActionHelp}, nil
	}

	if action == string(core.ActionList) {
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run can only be used when switching provider")
		}
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("Unexpected provider argument with --list")
		}
		return core.CliOptions{Action: core.ActionList}, nil
	}

	if action == string(core.ActionValidate) {
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run can only be used when switching provider")
		}
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("Unexpected provider argument with --validate")
		}
		return core.CliOptions{Action: core.ActionValidate}, nil
	}

	if len(positional) == 0 {
		return core.CliOptions{Action: core.ActionList}, nil
	}

	if len(positional) > 1 {
		return core.CliOptions{}, fmt.Errorf("Too many provider arguments: %s", strings.Join(positional, ", "))
	}

	return core.CliOptions{Action: core.ActionSwitch, ProviderInput: positional[0], DryRun: dryRun}, nil
}
