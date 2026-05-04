package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	services "github.com/HeminWon/proteus/internal/app"
	core "github.com/HeminWon/proteus/internal/cli"
)

func isHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h"
}

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
			return core.CliOptions{}, fmt.Errorf("unknown switch option: %s%s", arg, core.SuggestFlag(arg, []string{"--dry-run", "--help", "-h"}))
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

func parseValidateArgs(args []string) (core.CliOptions, error) {
	provider := ""
	concurrency := 5

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case isHelpFlag(arg):
			if len(args) > 1 {
				return core.CliOptions{}, fmt.Errorf("validate --help cannot be combined with other options")
			}
			return core.CliOptions{Action: core.ActionHelp, HelpCommand: string(core.ActionValidate)}, nil
		case arg == "--provider":
			if i+1 >= len(args) {
				return core.CliOptions{}, fmt.Errorf("missing value for --provider")
			}
			i++
			provider = args[i]
		case strings.HasPrefix(arg, "--provider="):
			return core.CliOptions{}, fmt.Errorf("unsupported validate option: %s (use `--provider <id>`)", arg)
		case arg == "--concurrency":
			if i+1 >= len(args) {
				return core.CliOptions{}, fmt.Errorf("missing value for --concurrency")
			}
			i++
			parsed, err := strconv.Atoi(args[i])
			if err != nil || parsed <= 0 {
				return core.CliOptions{}, fmt.Errorf("invalid --concurrency value: %s", args[i])
			}
			concurrency = parsed
		case strings.HasPrefix(arg, "--concurrency="):
			return core.CliOptions{}, fmt.Errorf("unsupported validate option: %s (use `--concurrency <n>`)", arg)
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, fmt.Errorf("unknown validate option: %s%s", arg, core.SuggestFlag(arg, []string{"--provider", "--concurrency", "--help", "-h"}))
		default:
			return core.CliOptions{}, fmt.Errorf("unexpected argument for validate: %s", arg)
		}
	}

	return core.CliOptions{Action: core.ActionValidate, ValidateProvider: provider, ValidateConcurrency: concurrency}, nil
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

func printHelpFor(command string) {
	switch command {
	case "", "global":
		printHelp()
	case string(core.ActionSwitch):
		fmt.Println("Usage:")
		fmt.Println("  proteus switch <provider-id|provider-name> [--dry-run]")
		fmt.Println()
		fmt.Println("Description:")
		fmt.Println("  Persist provider by overwriting ~/.claude/settings.json")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --dry-run        Preview switch plan without writing files")
		fmt.Println("  --help, -h       Show switch help")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  proteus switch anthropic --dry-run")
	case string(core.ActionLaunch):
		fmt.Println("Usage:")
		fmt.Println("  proteus launch <profile> [--dry-run]")
		fmt.Println("  proteus launch --list")
		fmt.Println()
		fmt.Println("Description:")
		fmt.Println("  Start claude with profile env in current process (no global file writes)")
		fmt.Println("  Note: launch writes profile-private settings, not global settings files")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --list           List launch profiles")
		fmt.Println("  --dry-run        Preview launch env and warnings")
		fmt.Println("  --help, -h       Show launch help")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  proteus launch default --dry-run")
	case string(core.ActionValidate):
		fmt.Println("Usage:")
		fmt.Println("  proteus validate [--provider <id>] [--concurrency <n>]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --provider       Validate only one provider")
		fmt.Println("  --concurrency    Live validation concurrency (default: 5)")
		fmt.Println("  --help, -h       Show validate help")
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  proteus list")
	fmt.Println("  proteus validate [--provider <id>] [--concurrency <n>]")
	fmt.Println("  proteus switch <provider-id|provider-name> [--dry-run]")
	fmt.Println("  proteus launch <profile> [--dry-run]")
	fmt.Println("  proteus launch --list")
	fmt.Println("  proteus --help")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list              List providers")
	fmt.Println("  validate          Validate providers.yaml and run live checks")
	fmt.Println("  switch            Persist provider by overwriting ~/.claude/settings.json")
	fmt.Println("  launch            Start claude with profile env in current process (no global file writes)")
	fmt.Println()
	fmt.Println("Tip:")
	fmt.Println("  Use `proteus <command> --help` for command-specific usage")
}

func run() error {
	parsed, err := parseArgs(os.Args[1:])
	if err != nil {
		return err
	}

	switch parsed.Action {
	case core.ActionHelp:
		printHelpFor(parsed.HelpCommand)
		return nil
	case core.ActionList:
		return services.ListProviders()
	case core.ActionValidate:
		return services.ValidateConfig(parsed.ValidateProvider, parsed.ValidateConcurrency)
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
