package main

import (
	"fmt"
	"os"

	services "github.com/HeminWon/proteus/internal/app"
	core "github.com/HeminWon/proteus/internal/cli"
)

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
