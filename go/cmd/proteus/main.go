package main

import (
	"fmt"
	"os"

	"github.com/HeminWon/proteus/go/internal/cli"
	"github.com/HeminWon/proteus/go/internal/core"
	"github.com/HeminWon/proteus/go/internal/services"
)

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
	parsed, err := cli.ParseArgs(os.Args[1:])
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
