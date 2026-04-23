package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HeminWon/proteus/internal/core"
	"github.com/HeminWon/proteus/internal/services"
)

func parseArgs(argv []string) (core.CliOptions, error) {
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
			return core.CliOptions{}, fmt.Errorf("unknown option: %s", arg)
		default:
			positional = append(positional, arg)
		}
	}

	switch core.CliAction(action) {
	case core.ActionHelp:
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run cannot be used with --help")
		}
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with --help")
		}
		return core.CliOptions{Action: core.ActionHelp}, nil
	case core.ActionList:
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run can only be used when switching provider")
		}
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with --list")
		}
		return core.CliOptions{Action: core.ActionList}, nil
	case core.ActionValidate:
		if dryRun {
			return core.CliOptions{}, fmt.Errorf("--dry-run can only be used when switching provider")
		}
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with --validate")
		}
		return core.CliOptions{Action: core.ActionValidate}, nil
	}

	if len(positional) == 0 {
		return core.CliOptions{Action: core.ActionList}, nil
	}
	if len(positional) > 1 {
		return core.CliOptions{}, fmt.Errorf("too many provider arguments: %s", strings.Join(positional, ", "))
	}
	return core.CliOptions{Action: core.ActionSwitch, ProviderInput: positional[0], DryRun: dryRun}, nil
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  proteus [provider-id|provider-name] [--dry-run]")
	fmt.Println("  proteus --list")
	fmt.Println("  proteus --validate")
	fmt.Println("  proteus --help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --list           List providers (default when no args)")
	fmt.Println("  --validate       Validate providers.yaml and run live checks")
	fmt.Println("  --dry-run        Preview switch plan without writing files")
	fmt.Println("  --help, -h       Show help")
}

func run() error {
	parsed, err := parseArgs(os.Args[1:])
	if err != nil {
		return err
	}

	switch parsed.Action {
	case core.ActionHelp:
		printHelp()
		return nil
	case core.ActionList:
		return services.ListProviders()
	case core.ActionValidate:
		return services.ValidateConfig()
	case core.ActionSwitch:
		return services.ApplyProvider(parsed.ProviderInput, parsed.DryRun)
	default:
		return fmt.Errorf("unsupported action: %s", parsed.Action)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		fmt.Fprintln(os.Stderr, "Tip: run with --help to see supported usage.")
		os.Exit(1)
	}
}
