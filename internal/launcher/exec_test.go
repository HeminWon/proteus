package launcher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HeminWon/proteus/internal/providers"
)

func TestResolveAuthTokenPriority(t *testing.T) {
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "auth-token")
	t.Setenv("ANTHROPIC_API_KEY", "api-key")
	ensurePathHasDummyClaude(t)

	config := providers.ProvidersConfig{
		Profiles: map[string]providers.Profile{"default": {Provider: "anthropic"}},
		Providers: []providers.Provider{{
			ID:   "anthropic",
			Name: "Anthropic",
			Claude: struct {
				Env    map[string]string `yaml:"env"`
				Models []string          `yaml:"models,omitempty"`
			}{Env: map[string]string{
				"ANTHROPIC_AUTH_TOKEN": "$ANTHROPIC_AUTH_TOKEN",
				"ANTHROPIC_API_KEY":    "$ANTHROPIC_API_KEY",
			}},
		}},
	}

	resolved, err := Resolve(config, "default")
	if err != nil {
		t.Fatalf("Resolve error = %v", err)
	}
	if resolved.TokenSource != "ANTHROPIC_AUTH_TOKEN" {
		t.Fatalf("TokenSource = %q, want ANTHROPIC_AUTH_TOKEN", resolved.TokenSource)
	}
	if len(resolved.CriticalWarns) != 0 {
		t.Fatalf("expected no critical warnings, got %v", resolved.CriticalWarns)
	}
}

func TestResolveApiKeyFallback(t *testing.T) {
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	t.Setenv("ANTHROPIC_API_KEY", "api-key")
	ensurePathHasDummyClaude(t)

	config := providers.ProvidersConfig{
		Profiles: map[string]providers.Profile{"default": {Provider: "anthropic"}},
		Providers: []providers.Provider{{
			ID:   "anthropic",
			Name: "Anthropic",
			Claude: struct {
				Env    map[string]string `yaml:"env"`
				Models []string          `yaml:"models,omitempty"`
			}{Env: map[string]string{
				"ANTHROPIC_AUTH_TOKEN": "$ANTHROPIC_AUTH_TOKEN",
				"ANTHROPIC_API_KEY":    "$ANTHROPIC_API_KEY",
			}},
		}},
	}

	resolved, err := Resolve(config, "default")
	if err != nil {
		t.Fatalf("Resolve error = %v", err)
	}
	if resolved.TokenSource != "ANTHROPIC_API_KEY" {
		t.Fatalf("TokenSource = %q, want ANTHROPIC_API_KEY", resolved.TokenSource)
	}
	if len(resolved.CriticalWarns) != 0 {
		t.Fatalf("expected no critical warnings, got %v", resolved.CriticalWarns)
	}
}

func TestResolveMissingAuthCritical(t *testing.T) {
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	ensurePathHasDummyClaude(t)

	config := providers.ProvidersConfig{
		Profiles: map[string]providers.Profile{"default": {Provider: "anthropic"}},
		Providers: []providers.Provider{{
			ID:   "anthropic",
			Name: "Anthropic",
			Claude: struct {
				Env    map[string]string `yaml:"env"`
				Models []string          `yaml:"models,omitempty"`
			}{Env: map[string]string{
				"ANTHROPIC_AUTH_TOKEN": "$ANTHROPIC_AUTH_TOKEN",
				"ANTHROPIC_API_KEY":    "$ANTHROPIC_API_KEY",
			}},
		}},
	}

	resolved, err := Resolve(config, "default")
	if err != nil {
		t.Fatalf("Resolve error = %v", err)
	}
	if resolved.TokenSource != "" {
		t.Fatalf("TokenSource = %q, want empty", resolved.TokenSource)
	}
	if len(resolved.CriticalWarns) == 0 {
		t.Fatalf("expected critical warning for missing auth")
	}
}

func ensurePathHasDummyClaude(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "claude")
	if err := os.WriteFile(file, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write dummy claude: %v", err)
	}
	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))
}
