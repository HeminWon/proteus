import type { Provider } from "../providers/types.js";
import type { ProvidersConfig } from "../providers/types.js";

export function validateProvidersConfigShape(parsed: unknown): ProvidersConfig {
  if (!parsed || typeof parsed !== "object") {
    throw new Error("providers.yaml is invalid: expected object root");
  }

  const config = parsed as ProvidersConfig;
  if (!Array.isArray(config.providers)) {
    throw new Error("providers.yaml is invalid: providers must be an array");
  }

  if (config.providers.length === 0) {
    throw new Error("providers.yaml is invalid: providers must not be empty");
  }

  const ids = new Set<string>();
  for (const provider of config.providers) {
    if (!provider?.id || typeof provider.id !== "string") {
      throw new Error("providers.yaml is invalid: each provider must have string id");
    }

    if (ids.has(provider.id)) {
      throw new Error(`providers.yaml is invalid: duplicate provider id '${provider.id}'`);
    }
    ids.add(provider.id);

    if (!provider.claude || typeof provider.claude !== "object") {
      throw new Error(`providers.yaml is invalid: provider '${provider.id}' missing claude`);
    }

    if (!provider.claude.env || typeof provider.claude.env !== "object") {
      throw new Error(
        `providers.yaml is invalid: provider '${provider.id}' missing claude.env`
      );
    }
  }

  return config;
}

export type { Provider };
