package services

import (
	"fmt"

	"github.com/HeminWon/proteus/internal/providers"
	"github.com/HeminWon/proteus/internal/store"
)

func getActiveProviderID(config providers.ProvidersConfig, cache store.CacheData) string {
	if cache.Active == nil || cache.Active.Claude == "" {
		return ""
	}
	for _, p := range config.Providers {
		if p.ID == cache.Active.Claude {
			return p.ID
		}
	}
	return ""
}

func ListProviders() error {
	loaded, err := providers.LoadProviders()
	if err != nil {
		return err
	}

	activeID := getActiveProviderID(loaded.Config, store.ReadCache())
	fmt.Printf("Config dir: %s\n\n", loaded.ConfigDir)
	fmt.Println("Available providers:")
	for _, p := range loaded.Config.Providers {
		active := ""
		if activeID != "" && p.ID == activeID {
			active = " ◀ active"
		}
		fmt.Printf("  %-16s %s%s\n", p.ID, p.Name, active)
	}

	return nil
}
