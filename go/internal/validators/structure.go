package validators

import (
	"fmt"

	"github.com/HeminWon/proteus/go/internal/providers"
)

func ValidateProvidersConfigShape(config providers.ProvidersConfig) error {
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

	return nil
}
