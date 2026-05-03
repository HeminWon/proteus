package main

import (
	"fmt"
	"os"
	"strings"

	services "github.com/HeminWon/proteus/internal/app"
	core "github.com/HeminWon/proteus/internal/cli"
)

func parseLaunchArgs(args []string) (core.CliOptions, error) {
	dryRun := false
	list := false
	positional := make([]string, 0)

	for _, arg := range args {
		switch {
		case arg == "--dry-run":
			dryRun = true
		case arg == "--list":
			list = true
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, fmt.Errorf("unknown launch option: %s", arg)
		default:
			positional = append(positional, arg)
		}
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

func parseSwitchArgs(args []string) (core.CliOptions, error) {
	dryRun := false
	positional := make([]string, 0)

	for _, arg := range args {
		switch {
		case arg == "--dry-run":
			dryRun = true
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, fmt.Errorf("unknown switch option: %s", arg)
		default:
			positional = append(positional, arg)
		}
	}

	if len(positional) == 0 {
		return core.CliOptions{}, fmt.Errorf("missing provider for switch")
	}
	if len(positional) > 1 {
		return core.CliOptions{}, fmt.Errorf("too many provider arguments: %s", strings.Join(positional, ", "))
	}

	return core.CliOptions{Action: core.ActionSwitch, ProviderInput: positional[0], DryRun: dryRun}, nil
}

func parseArgs(argv []string) (core.CliOptions, error) {
	action := ""
	positional := make([]string, 0)

	for _, arg := range argv {
		switch {
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
		case arg == "launch":
			if action != "" {
				return core.CliOptions{}, fmt.Errorf("launch cannot be combined with other actions")
			}
			return parseLaunchArgs(argv[1:])
		case arg == "switch":
			if action != "" {
				return core.CliOptions{}, fmt.Errorf("switch cannot be combined with other actions")
			}
			return parseSwitchArgs(argv[1:])
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, fmt.Errorf("unknown option: %s", arg)
		default:
			positional = append(positional, arg)
		}
	}

	switch core.CliAction(action) {
	case core.ActionHelp:
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with --help")
		}
		return core.CliOptions{Action: core.ActionHelp}, nil
	case core.ActionList:
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with --list")
		}
		return core.CliOptions{Action: core.ActionList}, nil
	case core.ActionValidate:
		if len(positional) > 0 {
			return core.CliOptions{}, fmt.Errorf("unexpected provider argument with --validate")
		}
		return core.CliOptions{Action: core.ActionValidate}, nil
	}

	if len(positional) == 0 {
		return core.CliOptions{Action: core.ActionHelp}, nil
	}

	return core.CliOptions{}, fmt.Errorf("unknown command: %s (run with --help to see supported commands)", positional[0])
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  proteus switch <provider-id|provider-name> [--dry-run]")
	fmt.Println("  proteus launch <profile> [--dry-run]")
	fmt.Println("  proteus launch --list")
	fmt.Println("  proteus --list")
	fmt.Println("  proteus --validate")
	fmt.Println("  proteus --help")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  switch            Persist provider by overwriting ~/.claude/settings.json")
	fmt.Println("  launch            Start claude with profile env in current process (no file writes)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --list           List providers")
	fmt.Println("  --validate       Validate providers.yaml and run live checks")
	fmt.Println("  --dry-run        Preview switch/launch plan without executing writes/exec")
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
	case core.ActionLaunch:
		return services.LaunchProfile(parsed.ProfileInput, parsed.DryRun, parsed.ListLaunch)
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
