package main

import (
	"fmt"

	core "github.com/HeminWon/proteus/internal/cli"
)

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
