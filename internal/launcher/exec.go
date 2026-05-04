package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	core "github.com/HeminWon/proteus/internal/cli"
	"github.com/HeminWon/proteus/internal/providers"
	store "github.com/HeminWon/proteus/internal/storage"
)

type ResolvedLaunch struct {
	Profile             string
	ProviderID          string
	ProviderName        string
	Runner              string
	RunnerPath          string
	RunnerArgs          []string
	ClaudeConfigDir     string
	PrivateSettingsPath string
	Provider            providers.Provider
	Env                 map[string]string
	ProviderEnvKeys     []string
	TokenSource         string
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
		profiles := availableProfiles(config)
		return ResolvedLaunch{}, fmt.Errorf("profile %q not found. Available profiles: %s%s", profile, strings.Join(profiles, ", "), core.SuggestProfile(profile, profiles))
	}

	provider := findProviderByID(config, binding.Provider)
	if provider == nil {
		return ResolvedLaunch{}, fmt.Errorf("profile %q references missing provider %q. Available providers: %s", profile, binding.Provider, strings.Join(availableProviders(config), ", "))
	}

	runner := strings.TrimSpace(binding.Runner)
	if runner == "" {
		runner = "claude"
	}
	runnerPath, err := exec.LookPath(runner)
	if err != nil {
		return ResolvedLaunch{}, fmt.Errorf("%s executable not found in PATH", runner)
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
			if k == "ANTHROPIC_AUTH_TOKEN" || k == "ANTHROPIC_API_KEY" {
				continue
			}
			warnings = append(warnings, fmt.Sprintf("WARN: env %s expanded to empty value", k))
		}
	}

	authToken := strings.TrimSpace(base["ANTHROPIC_AUTH_TOKEN"])
	apiKey := strings.TrimSpace(base["ANTHROPIC_API_KEY"])
	tokenSource := ""
	switch {
	case authToken != "":
		tokenSource = "ANTHROPIC_AUTH_TOKEN"
	case apiKey != "":
		tokenSource = "ANTHROPIC_API_KEY"
	default:
		critical = append(critical, "WARN[critical]: both ANTHROPIC_AUTH_TOKEN and ANTHROPIC_API_KEY are missing or empty; authentication will fail")
	}

	if authToken == "" && base["ANTHROPIC_AUTH_TOKEN"] != "" {
		warnings = append(warnings, "WARN: ANTHROPIC_AUTH_TOKEN resolves to whitespace-only value")
	}
	if apiKey == "" && base["ANTHROPIC_API_KEY"] != "" {
		warnings = append(warnings, "WARN: ANTHROPIC_API_KEY resolves to whitespace-only value")
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
		Runner:              runner,
		RunnerPath:          runnerPath,
		RunnerArgs:          append([]string{}, binding.Args...),
		ClaudeConfigDir:     claudeConfigDir,
		PrivateSettingsPath: privateSettingsPath,
		Provider:            *provider,
		Env:                 base,
		ProviderEnvKeys:     providerKeys,
		TokenSource:         tokenSource,
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
	argv := append([]string{resolved.Runner}, resolved.RunnerArgs...)
	return syscall.Exec(resolved.RunnerPath, argv, env)
}
