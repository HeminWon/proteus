package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type appConfig struct {
	ConfigDir string `json:"config_dir"`
}

func readAppConfig() appConfig {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}

	configPath := filepath.Join(xdgConfigHome, "proteus", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return appConfig{}
	}

	var parsed appConfig
	if err := json.Unmarshal(data, &parsed); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", configPath, err)
		return appConfig{}
	}

	return parsed
}

func ResolveConfigDir() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	xdgProteusDir := filepath.Join(xdgConfigHome, "proteus")

	appCfg := readAppConfig()
	type candidate struct {
		dir   string
		label string
	}
	candidates := make([]candidate, 0, 2)
	if appCfg.ConfigDir != "" {
		expanded := appCfg.ConfigDir
		if strings.HasPrefix(expanded, "~") {
			expanded = strings.Replace(expanded, "~", os.Getenv("HOME"), 1)
		}
		candidates = append(candidates, candidate{dir: expanded, label: "config.json (config_dir)"})
	}
	candidates = append(candidates, candidate{dir: xdgProteusDir, label: "XDG (" + xdgProteusDir + ")"})

	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c.dir, "providers.yaml")); err == nil {
			return c.dir, nil
		}
	}

	searched := make([]string, 0, len(candidates))
	for _, c := range candidates {
		searched = append(searched, "  - "+filepath.Join(c.dir, "providers.yaml"))
	}
	configJSONPath := filepath.Join(xdgProteusDir, "config.json")

	return "", fmt.Errorf("providers.yaml not found.\nSearched:\n%s\n\nPlace your providers.yaml at:\n  %s\nOr set config_dir in %s", strings.Join(searched, "\n"), filepath.Join(xdgProteusDir, "providers.yaml"), configJSONPath)
}

func validateProvidersConfigShape(config ProvidersConfig) error {
	if len(config.Providers) == 0 {
		return fmt.Errorf("providers.yaml is invalid: providers must not be empty")
	}

	ids := map[string]struct{}{}
	for _, provider := range config.Providers {
		if provider.ID == "" {
			return fmt.Errorf("providers.yaml is invalid: each provider must have string id")
		}

		if _, exists := ids[provider.ID]; exists {
			return fmt.Errorf("providers.yaml is invalid: duplicate provider id '%s'", provider.ID)
		}
		ids[provider.ID] = struct{}{}

		if provider.Claude.Env == nil {
			return fmt.Errorf("providers.yaml is invalid: provider '%s' missing claude.env", provider.ID)
		}
	}

	for profile, binding := range config.Profiles {
		if strings.TrimSpace(profile) == "" {
			return fmt.Errorf("providers.yaml is invalid: profile name must not be empty")
		}
		if strings.TrimSpace(binding.Provider) == "" {
			return fmt.Errorf("providers.yaml is invalid: profile '%s' missing provider", profile)
		}
	}

	return nil
}

func LoadProviders() (LoadProvidersResult, error) {
	configDir, err := ResolveConfigDir()
	if err != nil {
		return LoadProvidersResult{}, err
	}

	raw, err := os.ReadFile(filepath.Join(configDir, "providers.yaml"))
	if err != nil {
		return LoadProvidersResult{}, err
	}

	var config ProvidersConfig
	if err := yaml.Unmarshal(raw, &config); err != nil {
		return LoadProvidersResult{}, err
	}

	if err := validateProvidersConfigShape(config); err != nil {
		return LoadProvidersResult{}, err
	}

	return LoadProvidersResult{Config: config, ConfigDir: configDir}, nil
}
