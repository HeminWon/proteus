package services

import (
	"fmt"

	"github.com/HeminWon/proteus/internal/cli"
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

func validateProvidersLive(providersList []providers.Provider, concurrency int) []validators.LiveValidationResult {
	results := make([]validators.LiveValidationResult, 0, len(providersList))
	for i := 0; i < len(providersList); i += concurrency {
		end := i + concurrency
		if end > len(providersList) {
			end = len(providersList)
		}
		batch := providersList[i:end]
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

func ValidateConfig() error {
	loaded, err := providers.LoadProviders()
	if err != nil {
		return err
	}

	cache := store.ReadCache()
	activeProviderID := getActiveProviderID(loaded.Config, cache)
	if activeProviderID == "" {
		activeProviderID = "unset"
	}

	fmt.Println("providers.yaml is valid.")
	fmt.Printf("- config dir: %s\n", loaded.ConfigDir)
	fmt.Printf("- version: %d\n", loaded.Config.Version)
	fmt.Printf("- active (cache): %s\n", activeProviderID)
	fmt.Printf("- providers: %d\n", len(loaded.Config.Providers))

	concurrency := 5
	fmt.Printf("- live validation: enabled (HTTP endpoint, concurrency=%d)\n", concurrency)

	results := validateProvidersLive(loaded.Config.Providers, concurrency)
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
