export type JsonObject = Record<string, unknown>;

export interface SettingsReadResult {
  exists: boolean;
  data: JsonObject;
}

export interface SwitchPlan {
  fromProvider: string | null;
  toProvider: string;
  envAdded: string[];
  envUpdated: string[];
  envRemoved: string[];
  availableModelsChanged: boolean;
  backupRequired: boolean;
}

export interface CacheData {
  active?: {
    claude?: string;
  };
}

export interface CliOptions {
  action: "list" | "validate" | "switch" | "help";
  providerInput?: string;
  dryRun: boolean;
}
