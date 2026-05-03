package services

import (
	"fmt"
	"strings"

	core "github.com/HeminWon/proteus/internal/cli"
	"github.com/HeminWon/proteus/internal/providers"
	"github.com/HeminWon/proteus/internal/storage"
	"github.com/HeminWon/proteus/internal/term"
	"github.com/HeminWon/proteus/internal/validators"
)

func formatLatency(latencyMs *int64, status string) string {
	if latencyMs == nil {
		return term.Colorize("n/a", "33")
	}
	raw := fmt.Sprintf("%dms", *latencyMs)
	if *latencyMs >= int64(core.HighLatencyMs) {
		return term.Colorize(raw, "33")
	}
	if status == "ok" {
		return term.Colorize(raw, "32")
	}
	if status == "fail" {
		return term.Colorize(raw, "31")
	}
	return raw
}

func validateProvidersLive(providersList []providers.Provider, providerFilter string, concurrency int) []validators.LiveValidationResult {
	selected := providersList
	if providerFilter != "" {
		filtered := make([]providers.Provider, 0, len(providersList))
		for _, p := range providersList {
			if p.ID == providerFilter {
				filtered = append(filtered, p)
				break
			}
		}
		selected = filtered
	}

	results := make([]validators.LiveValidationResult, 0, len(selected))
	for i := 0; i < len(selected); i += concurrency {
		end := i + concurrency
		if end > len(selected) {
			end = len(selected)
		}
		batch := selected[i:end]
		batchResults := make(chan validators.LiveValidationResult, len(batch))
		for _, provider := range batch {
			p := provider
			go func() {
				batchResults <- validators.ValidateProviderLive(p)
			}()
		}
		for range batch {
			results = append(results, <-batchResults)
		}
		close(batchResults)
	}
	return results
}

func ValidateConfig(providerFilter string, concurrency int) error {
	loaded, err := providers.LoadProviders()
	if err != nil {
		return err
	}

	cache := store.ReadCache()
	activeProviderID := getActiveProviderID(loaded.Config, cache)
	if activeProviderID == "" {
		activeProviderID = "unset"
	}

	if concurrency <= 0 {
		concurrency = 5
	}

	if providerFilter != "" {
		exists := false
		available := make([]string, 0, len(loaded.Config.Providers))
		for _, p := range loaded.Config.Providers {
			available = append(available, p.ID)
			if p.ID == providerFilter {
				exists = true
			}
		}
		if !exists {
			return fmt.Errorf("provider %q not found. Available: %s%s", providerFilter, strings.Join(available, ", "), core.SuggestProvider(providerFilter, available))
		}
	}

	fmt.Println("providers.yaml is valid.")
	fmt.Printf("- config dir: %s\n", loaded.ConfigDir)
	fmt.Printf("- version: %d\n", loaded.Config.Version)
	fmt.Printf("- active (cache): %s\n", activeProviderID)
	fmt.Printf("- providers: %d\n", len(loaded.Config.Providers))
	if providerFilter == "" {
		fmt.Printf("- validate target: all\n")
	} else {
		fmt.Printf("- validate target: %s\n", providerFilter)
	}
	fmt.Printf("- live validation: enabled (HTTP endpoint, concurrency=%d)\n", concurrency)

	results := validateProvidersLive(loaded.Config.Providers, providerFilter, concurrency)
	failed := 0
	for _, result := range results {
		mark := "FAIL"
		switch result.Status {
		case "ok":
			mark = "OK"
		case "skip":
			mark = "SKIP"
		}
		markDisplay := term.ColorStatus(result.Status, mark)
		latencyDisplay := formatLatency(result.LatencyMs, result.Status)
		fmt.Printf("  [%s] %s: %s | latency=%s\n", markDisplay, result.ProviderID, result.Detail, latencyDisplay)
		if result.Status == "fail" {
			failed++
		}
	}

	if failed > 0 {
		return fmt.Errorf("live validation failed for %d provider(s)", failed)
	}
	return nil
}
