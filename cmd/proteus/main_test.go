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

	for _, cmd := range canonicalCommands {
		got, err = parseArgs([]string{"--help", cmd})
		if err != nil {
			t.Fatalf("parseArgs error = %v", err)
		}
		if got.HelpCommand != cmd {
			t.Fatalf("help command = %q, want %q", got.HelpCommand, cmd)
		}
	}

	_, err = parseArgs([]string{"help"})
	if err == nil {
		t.Fatalf("expected error for removed help command")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "unknown command: help") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseValidateSupportsEqualsForms(t *testing.T) {
	got, err := parseArgs([]string{"validate", "--provider=anthropic", "--concurrency=3"})
	if err != nil {
		t.Fatalf("parseArgs error = %v", err)
	}
	if got.ValidateProvider != "anthropic" {
		t.Fatalf("provider = %q, want anthropic", got.ValidateProvider)
	}
	if got.ValidateConcurrency != 3 {
		t.Fatalf("concurrency = %d, want 3", got.ValidateConcurrency)
	}

	got, err = parseArgs([]string{"validate", "--provider", "anthropic", "--concurrency=4"})
	if err != nil {
		t.Fatalf("parseArgs mixed form error = %v", err)
	}
	if got.ValidateProvider != "anthropic" || got.ValidateConcurrency != 4 {
		t.Fatalf("unexpected mixed parse: %+v", got)
	}

	got, err = parseArgs([]string{"validate", "--provider=anthropic", "--concurrency", "5"})
	if err != nil {
		t.Fatalf("parseArgs mixed form error = %v", err)
	}
	if got.ValidateProvider != "anthropic" || got.ValidateConcurrency != 5 {
		t.Fatalf("unexpected mixed parse: %+v", got)
	}
}

func TestParseValidateEqualsFormErrors(t *testing.T) {
	_, err := parseArgs([]string{"validate", "--provider="})
	if err == nil {
		t.Fatalf("expected error for empty --provider=")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "missing value for --provider") {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = parseArgs([]string{"validate", "--concurrency="})
	if err == nil {
		t.Fatalf("expected error for empty --concurrency=")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "missing value for --concurrency") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseValidateConcurrencyValidation(t *testing.T) {
	for _, value := range []string{"0", "-1", "abc"} {
		_, err := parseArgs([]string{"validate", "--concurrency=" + value})
		if err == nil {
			t.Fatalf("expected invalid concurrency error for value %q", value)
		}
		if got := err.Error(); got == "" || !strings.Contains(got, "invalid --concurrency value") {
			t.Fatalf("unexpected error for value %q: %v", value, err)
		}
	}
}

func TestParseSwitchArgumentRules(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantErrContains string
	}{
		{name: "missing provider", args: []string{"switch"}, wantErrContains: "missing provider for switch"},
		{name: "too many provider", args: []string{"switch", "a", "b"}, wantErrContains: "too many provider arguments"},
		{name: "help with provider", args: []string{"switch", "--help", "anthropic"}, wantErrContains: "unexpected provider argument with switch --help"},
		{name: "help with dry-run", args: []string{"switch", "--help", "--dry-run"}, wantErrContains: "switch --help cannot be combined with other options"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseArgs(tt.args)
			if err == nil {
				t.Fatalf("expected error")
			}
			if got := err.Error(); !strings.Contains(got, tt.wantErrContains) {
				t.Fatalf("error = %q, want contains %q", got, tt.wantErrContains)
			}
		})
	}
}

func TestParseLaunchArgumentRules(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantErrContains string
	}{
		{name: "missing profile", args: []string{"launch"}, wantErrContains: "missing profile for launch"},
		{name: "too many profile", args: []string{"launch", "a", "b"}, wantErrContains: "too many profile arguments"},
		{name: "help with profile", args: []string{"launch", "--help", "default"}, wantErrContains: "unexpected profile argument with launch --help"},
		{name: "help with list", args: []string{"launch", "--help", "--list"}, wantErrContains: "launch --help cannot be combined with other options"},
		{name: "list with dry-run", args: []string{"launch", "--list", "--dry-run"}, wantErrContains: "--dry-run cannot be used with launch --list"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseArgs(tt.args)
			if err == nil {
				t.Fatalf("expected error")
			}
			if got := err.Error(); !strings.Contains(got, tt.wantErrContains) {
				t.Fatalf("error = %q, want contains %q", got, tt.wantErrContains)
			}
		})
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
