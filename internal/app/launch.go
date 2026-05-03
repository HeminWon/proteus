package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/HeminWon/proteus/internal/launcher"
	"github.com/HeminWon/proteus/internal/providers"
	store "github.com/HeminWon/proteus/internal/storage"
)

func maskEnvValue(key, value string) string {
	upper := strings.ToUpper(key)
	if strings.Contains(upper, "KEY") || strings.Contains(upper, "TOKEN") || strings.Contains(upper, "SECRET") {
		if value == "" {
			return ""
		}
		if len(value) <= 6 {
			return "***"
		}
		return value[:3] + "***...***"
	}
	return value
}

func printResolvedDryRun(resolved launcher.ResolvedLaunch) {
	fmt.Printf("Profile:  %s\n", resolved.Profile)
	fmt.Printf("Provider: %s (%s)\n", resolved.ProviderID, resolved.ProviderName)
	fmt.Printf("Private settings: %s\n", resolved.PrivateSettingsPath)
	fmt.Printf("CLAUDE_CONFIG_DIR: %s\n", resolved.ClaudeConfigDir)
	if resolved.TokenSource == "" {
		fmt.Printf("Auth source: %s\n\n", "(missing)")
	} else {
		fmt.Printf("Auth source: %s\n\n", resolved.TokenSource)
	}
	fmt.Println("Env:")
	for _, key := range resolved.ProviderEnvKeys {
		fmt.Printf("  %-20s = %s\n", key, maskEnvValue(key, resolved.Env[key]))
	}
	for _, w := range resolved.Warnings {
		fmt.Println(w)
	}
	for _, w := range resolved.CriticalWarns {
		fmt.Println(w)
	}
}

func listLaunchProfiles(config providers.ProvidersConfig) {
	if len(config.Profiles) == 0 {
		fmt.Println("No profiles configured.")
		return
	}

	providerIDs := map[string]struct{}{}
	for _, p := range config.Providers {
		providerIDs[p.ID] = struct{}{}
	}

	profiles := make([]string, 0, len(config.Profiles))
	for profile := range config.Profiles {
		profiles = append(profiles, profile)
	}
	sort.Strings(profiles)

	fmt.Println("Launch profiles:")
	for _, profile := range profiles {
		binding := config.Profiles[profile]
		status := "ok"
		if _, exists := providerIDs[binding.Provider]; !exists {
			status = "missing-provider"
		}
		fmt.Printf("  %-16s provider=%-16s status=%s\n", profile, binding.Provider, status)
	}
}

func LaunchProfile(profile string, dryRun bool, list bool) error {
	loaded, err := providers.LoadProviders()
	if err != nil {
		return err
	}

	if list {
		listLaunchProfiles(loaded.Config)
		return nil
	}

	resolved, err := launcher.Resolve(loaded.Config, profile)
	if err != nil {
		return err
	}

	settings, err := store.ReadSettings()
	if err != nil {
		return store.WrapSettingsParseError(err)
	}
	nextSettings := BuildNextSettings(settings.Data, resolved.Provider)

	if dryRun {
		printResolvedDryRun(resolved)
		return nil
	}

	if err := store.WriteSettingsAt(resolved.PrivateSettingsPath, nextSettings); err != nil {
		return err
	}

	for _, w := range resolved.Warnings {
		fmt.Println(w)
	}
	for _, w := range resolved.CriticalWarns {
		fmt.Println(w)
	}

	return launcher.ExecResolved(resolved)
}
