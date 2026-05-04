package main

import (
	"fmt"
	"strconv"
	"strings"

	core "github.com/HeminWon/proteus/internal/cli"
)

var validateFlags = []string{"--provider", "--concurrency", "--help", "-h"}

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
			provider = strings.TrimPrefix(arg, "--provider=")
			if provider == "" {
				return core.CliOptions{}, fmt.Errorf("missing value for --provider")
			}
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
			value := strings.TrimPrefix(arg, "--concurrency=")
			if value == "" {
				return core.CliOptions{}, fmt.Errorf("missing value for --concurrency")
			}
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed <= 0 {
				return core.CliOptions{}, fmt.Errorf("invalid --concurrency value: %s", value)
			}
			concurrency = parsed
		case strings.HasPrefix(arg, "-"):
			return core.CliOptions{}, unknownOptionError(string(core.ActionValidate), arg, validateFlags)
		default:
			return core.CliOptions{}, fmt.Errorf("unexpected argument for validate: %s", arg)
		}
	}

	return core.CliOptions{Action: core.ActionValidate, ValidateProvider: provider, ValidateConcurrency: concurrency}, nil
}
