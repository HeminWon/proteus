package services

import (
	"fmt"
	"sort"

	"github.com/HeminWon/proteus/go/internal/providers"
	"github.com/HeminWon/proteus/go/internal/store"
)

type SwitchPlan struct {
	FromProvider           string
	ToProvider             string
	EnvAdded               []string
	EnvUpdated             []string
	EnvRemoved             []string
	AvailableModelsChanged bool
	BackupRequired         bool
}

func findProviderByID(config providers.ProvidersConfig, id string) *providers.Provider {
	for i := range config.Providers {
		if config.Providers[i].ID == id {
			return &config.Providers[i]
		}
	}
	return nil
}

func findProviderByInput(config providers.ProvidersConfig, input string) *providers.Provider {
	if byID := findProviderByID(config, input); byID != nil {
		return byID
	}
	for i := range config.Providers {
		if config.Providers[i].Name == input {
			return &config.Providers[i]
		}
	}
	return nil
}

func buildNextSettings(current store.JsonObject, provider providers.Provider) store.JsonObject {
	next := store.JsonObject{}
	for k, v := range current {
		next[k] = v
	}

	env := map[string]any{}
	for k, v := range provider.Claude.Env {
		env[k] = v
	}
	next["env"] = env

	if len(provider.Claude.Models) > 0 {
		models := make([]any, 0, len(provider.Claude.Models))
		for _, m := range provider.Claude.Models {
			models = append(models, m)
		}
		next["availableModels"] = models
	} else {
		delete(next, "availableModels")
	}

	return next
}

func asStringMap(value any) map[string]string {
	obj, ok := value.(map[string]any)
	if !ok {
		return map[string]string{}
	}
	result := map[string]string{}
	for k, v := range obj {
		s, ok := v.(string)
		if ok {
			result[k] = s
		}
	}
	return result
}

func buildSwitchPlan(activeProviderID string, currentSettings store.JsonObject, nextSettings store.JsonObject, targetProvider string, backupRequired bool) SwitchPlan {
	beforeEnv := asStringMap(currentSettings["env"])
	afterEnv := asStringMap(nextSettings["env"])

	envAdded := make([]string, 0)
	envUpdated := make([]string, 0)
	envRemoved := make([]string, 0)

	for key, value := range afterEnv {
		if before, exists := beforeEnv[key]; !exists {
			envAdded = append(envAdded, key)
		} else if before != value {
			envUpdated = append(envUpdated, key)
		}
	}

	for key := range beforeEnv {
		if _, exists := afterEnv[key]; !exists {
			envRemoved = append(envRemoved, key)
		}
	}

	sort.Strings(envAdded)
	sort.Strings(envUpdated)
	sort.Strings(envRemoved)

	from := activeProviderID
	if from == "" {
		from = "(unset)"
	}

	_, beforeHasModels := currentSettings["availableModels"]
	_, afterHasModels := nextSettings["availableModels"]
	modelsChanged := beforeHasModels != afterHasModels
	if beforeHasModels && afterHasModels {
		modelsChanged = fmt.Sprintf("%v", currentSettings["availableModels"]) != fmt.Sprintf("%v", nextSettings["availableModels"])
	}

	return SwitchPlan{
		FromProvider:           from,
		ToProvider:             targetProvider,
		EnvAdded:               envAdded,
		EnvUpdated:             envUpdated,
		EnvRemoved:             envRemoved,
		AvailableModelsChanged: modelsChanged,
		BackupRequired:         backupRequired,
	}
}

func printSwitchPlan(plan SwitchPlan) {
	fmt.Printf("Plan: %s -> %s\n", plan.FromProvider, plan.ToProvider)
	fmt.Println("- mode: overwrite-env")
	fmt.Printf("- settings: %s\n", store.SettingsPath())
	fmt.Printf("- cache: %s\n", store.CachePath())
	if plan.BackupRequired {
		fmt.Println("- backup: create")
	} else {
		fmt.Println("- backup: skip (file missing)")
	}

	printKeys := func(label string, keys []string) {
		if len(keys) == 0 {
			fmt.Printf("- %s: none\n", label)
			return
		}
		fmt.Printf("- %s: %s\n", label, join(keys, ", "))
	}
	printKeys("env added", plan.EnvAdded)
	printKeys("env updated", plan.EnvUpdated)
	printKeys("env removed", plan.EnvRemoved)
	if plan.AvailableModelsChanged {
		fmt.Println("- availableModels: changed")
	} else {
		fmt.Println("- availableModels: unchanged")
	}
}

func join(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	out := items[0]
	for i := 1; i < len(items); i++ {
		out += sep + items[i]
	}
	return out
}

func ApplyProvider(input string, dryRun bool) error {
	loaded, err := providers.LoadProviders()
	if err != nil {
		return err
	}

	cache := store.ReadCache()
	activeProviderID := getActiveProviderID(loaded.Config, cache)
	provider := findProviderByInput(loaded.Config, input)
	if provider == nil {
		available := make([]string, 0, len(loaded.Config.Providers))
		for _, p := range loaded.Config.Providers {
			available = append(available, p.ID)
		}
		return fmt.Errorf("Provider \"%s\" not found. Available: %s", input, join(available, ", "))
	}

	settings, err := store.ReadSettings()
	if err != nil {
		return store.WrapSettingsParseError(err)
	}

	nextSettings := buildNextSettings(settings.Data, *provider)
	plan := buildSwitchPlan(activeProviderID, settings.Data, nextSettings, provider.ID, settings.Exists)

	if dryRun {
		printSwitchPlan(plan)
		return nil
	}

	backupPath, err := store.CreateBackupIfNeeded(settings.Exists)
	if err != nil {
		return err
	}

	if err := store.WriteSettings(nextSettings); err != nil {
		if backupPath != "" {
			_ = store.RestoreFromBackup(backupPath)
		}
		return err
	}

	nextCache := store.CacheData{Active: &struct {
		Claude string `json:"claude,omitempty"`
	}{Claude: provider.ID}}
	if err := store.WriteCache(nextCache); err != nil {
		if backupPath != "" {
			_ = store.RestoreFromBackup(backupPath)
		}
		return err
	}

	fmt.Printf("Switched to: %s (%s)\n", provider.Name, provider.ID)
	if backupPath != "" {
		fmt.Printf("Backup: %s\n", backupPath)
	}
	fmt.Println("Mode: overwrite-env")
	return nil
}
