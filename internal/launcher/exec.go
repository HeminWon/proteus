package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"github.com/HeminWon/proteus/internal/providers"
	store "github.com/HeminWon/proteus/internal/storage"
)

type ResolvedLaunch struct {
	Profile             string
	ProviderID          string
	ProviderName        string
	ClaudePath          string
	ClaudeConfigDir     string
	PrivateSettingsPath string
	Provider            providers.Provider
	Env                 map[string]string
	ProviderEnvKeys     []string
	Warnings            []string
	CriticalWarns       []string
}

func findProviderByID(config providers.ProvidersConfig, id string) *providers.Provider {
	for i := range config.Providers {
		if config.Providers[i].ID == id {
			return &config.Providers[i]
		}
	}
	return nil
}

func availableProfiles(config providers.ProvidersConfig) []string {
	keys := make([]string, 0, len(config.Profiles))
	for name := range config.Profiles {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return keys
}

func availableProviders(config providers.ProvidersConfig) []string {
	keys := make([]string, 0, len(config.Providers))
	for _, p := range config.Providers {
		keys = append(keys, p.ID)
	}
	sort.Strings(keys)
	return keys
}

func Resolve(config providers.ProvidersConfig, profile string) (ResolvedLaunch, error) {
	binding, ok := config.Profiles[profile]
	if !ok {
		return ResolvedLaunch{}, fmt.Errorf("profile %q not found. Available profiles: %s", profile, strings.Join(availableProfiles(config), ", "))
	}

	provider := findProviderByID(config, binding.Provider)
	if provider == nil {
		return ResolvedLaunch{}, fmt.Errorf("profile %q references missing provider %q. Available providers: %s", profile, binding.Provider, strings.Join(availableProviders(config), ", "))
	}

	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return ResolvedLaunch{}, fmt.Errorf("claude executable not found in PATH")
	}

	base := environToMap(os.Environ())
	for _, k := range []string{"ANTHROPIC_MODEL", "ANTHROPIC_DEFAULT_SONNET_MODEL", "ANTHROPIC_DEFAULT_OPUS_MODEL", "ANTHROPIC_DEFAULT_HAIKU_MODEL"} {
		delete(base, k)
	}
	providerKeys := make([]string, 0, len(provider.Claude.Env))
	warnings := make([]string, 0)
	critical := make([]string, 0)

	for k, v := range provider.Claude.Env {
		expanded := os.ExpandEnv(v)
		base[k] = expanded
		providerKeys = append(providerKeys, k)
		if expanded == "" {
			warnings = append(warnings, fmt.Sprintf("WARN: env %s expanded to empty value", k))
		}
		if k == "ANTHROPIC_API_KEY" && expanded == "" {
			critical = append(critical, "WARN[critical]: ANTHROPIC_API_KEY is empty and auth may fail")
		}
	}

	if len(providerKeys) == 0 {
		return ResolvedLaunch{}, fmt.Errorf("provider %q has no claude.env entries; launch configuration is empty", provider.ID)
	}

	claudeConfigDir := store.LaunchProfileConfigDir(profile)
	privateSettingsPath := store.LaunchProfileSettingsPath(profile)
	base["CLAUDE_CONFIG_DIR"] = claudeConfigDir

	sort.Strings(providerKeys)
	return ResolvedLaunch{
		Profile:             profile,
		ProviderID:          provider.ID,
		ProviderName:        provider.Name,
		ClaudePath:          claudePath,
		ClaudeConfigDir:     claudeConfigDir,
		PrivateSettingsPath: privateSettingsPath,
		Provider:            *provider,
		Env:                 base,
		ProviderEnvKeys:     providerKeys,
		Warnings:            warnings,
		CriticalWarns:       critical,
	}, nil
}

func environToMap(values []string) map[string]string {
	result := make(map[string]string, len(values))
	for _, item := range values {
		idx := strings.Index(item, "=")
		if idx <= 0 {
			continue
		}
		result[item[:idx]] = item[idx+1:]
	}
	return result
}

func mapToEnviron(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+values[k])
	}
	return out
}

func ExecResolved(resolved ResolvedLaunch) error {
	env := mapToEnviron(resolved.Env)
	return syscall.Exec(resolved.ClaudePath, []string{"claude"}, env)
}
