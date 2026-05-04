package main

import (
	"fmt"

	core "github.com/HeminWon/proteus/internal/cli"
)

type helpDoc struct {
	usage       []string
	description []string
	options     []string
	example     []string
}

var commandDescriptions = map[string]string{
	commandList:                 "List providers",
	string(core.ActionValidate): "Validate providers.yaml and run live checks",
	string(core.ActionSwitch):   "Persist provider by overwriting ~/.claude/settings.json",
	string(core.ActionLaunch):   "Start claude with profile env in current process (no global file writes)",
}

var helpDocs = map[string]helpDoc{
	string(core.ActionSwitch): {
		usage:       []string{"proteus switch <provider-id|provider-name> [--dry-run]"},
		description: []string{"Persist provider by overwriting ~/.claude/settings.json"},
		options: []string{
			"--dry-run        Preview switch plan without writing files",
			"--help, -h       Show switch help",
		},
		example: []string{"proteus switch anthropic --dry-run"},
	},
	string(core.ActionLaunch): {
		usage: []string{"proteus launch <profile> [--dry-run]", "proteus launch --list"},
		description: []string{
			"Start claude with profile env in current process (no global file writes)",
			"Note: launch writes profile-private settings, not global settings files",
		},
		options: []string{
			"--list           List launch profiles",
			"--dry-run        Preview launch env and warnings",
			"--help, -h       Show launch help",
		},
		example: []string{"proteus launch default --dry-run"},
	},
	string(core.ActionValidate): {
		usage: []string{"proteus validate [--provider <id>] [--concurrency <n>]"},
		options: []string{
			"--provider       Validate only one provider",
			"--concurrency    Live validation concurrency (default: 5)",
			"--help, -h       Show validate help",
		},
	},
}

func renderSection(title string, lines []string, indent bool) {
	if len(lines) == 0 {
		return
	}
	fmt.Println(title + ":")
	for _, line := range lines {
		if indent {
			fmt.Println("  " + line)
			continue
		}
		fmt.Println(line)
	}
	fmt.Println()
}

func printCommandHelp(command string) {
	doc, ok := helpDocs[command]
	if !ok {
		printHelp()
		return
	}

	renderSection("Usage", doc.usage, true)
	renderSection("Description", doc.description, true)
	renderSection("Options", doc.options, true)
	if len(doc.example) > 0 {
		renderSection("Example", doc.example, true)
	}
}

func printHelpFor(command string) {
	switch command {
	case "", globalHelpAlias:
		printHelp()
	case commandList, string(core.ActionSwitch), string(core.ActionLaunch), string(core.ActionValidate):
		printCommandHelp(command)
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
	for _, cmd := range canonicalCommands {
		fmt.Printf("  %-17s %s\n", cmd, commandDescriptions[cmd])
	}
	fmt.Println()
	fmt.Println("Tip:")
	fmt.Println("  Use `proteus <command> --help` for command-specific usage")
}
