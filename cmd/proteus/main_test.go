package main

import (
	"strings"
	"testing"
)

func TestParseArgsCanonicalCommands(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantAction string
	}{
		{name: "list", args: []string{"list"}, wantAction: "list"},
		{name: "validate", args: []string{"validate"}, wantAction: "validate"},
		{name: "switch", args: []string{"switch", "anthropic"}, wantAction: "switch"},
		{name: "launch", args: []string{"launch", "default"}, wantAction: "launch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if err != nil {
				t.Fatalf("parseArgs error = %v", err)
			}
			if string(got.Action) != tt.wantAction {
				t.Fatalf("action = %q, want %q", got.Action, tt.wantAction)
			}
		})
	}
}

func TestParseArgsLegacyCompatibility(t *testing.T) {
	_, err := parseArgs([]string{"--list"})
	if err == nil {
		t.Fatalf("expected error for removed legacy --list")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "unknown option: --list") {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = parseArgs([]string{"--validate", "--concurrency", "3"})
	if err == nil {
		t.Fatalf("expected error for removed legacy --validate")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "unknown option: --validate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseArgsHelpDispatch(t *testing.T) {
	got, err := parseArgs([]string{"--help"})
	if err != nil {
		t.Fatalf("parseArgs error = %v", err)
	}
	if string(got.Action) != "help" || got.HelpCommand != "" {
		t.Fatalf("unexpected help parse result: %+v", got)
	}

	got, err = parseArgs([]string{"switch", "--help"})
	if err != nil {
		t.Fatalf("parseArgs error = %v", err)
	}
	if got.HelpCommand != "switch" {
		t.Fatalf("help command = %q, want switch", got.HelpCommand)
	}

	got, err = parseArgs([]string{"launch", "--help"})
	if err != nil {
		t.Fatalf("parseArgs error = %v", err)
	}
	if got.HelpCommand != "launch" {
		t.Fatalf("help command = %q, want launch", got.HelpCommand)
	}

	_, err = parseArgs([]string{"help"})
	if err == nil {
		t.Fatalf("expected error for removed help command")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "unknown command: help") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseValidateProviderRequiresSpaceSeparatedValue(t *testing.T) {
	_, err := parseArgs([]string{"validate", "--provider=anthropic"})
	if err == nil {
		t.Fatalf("expected error for unsupported --provider=<id> form")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "unsupported validate option: --provider=anthropic") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseArgsSuggestions(t *testing.T) {
	_, err := parseArgs([]string{"swith"})
	if err == nil {
		t.Fatalf("expected error for unknown command")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "Did you mean `switch`") {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = parseArgs([]string{"switch", "--dryun"})
	if err == nil {
		t.Fatalf("expected error for unknown flag")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "Did you mean `--dry-run`") {
		t.Fatalf("unexpected error: %v", err)
	}
}
