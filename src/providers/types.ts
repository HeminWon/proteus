export interface Provider {
  id: string;
  name: string;
  claude: {
    env: Record<string, string>;
    models?: string[];
  };
  codex?: {
    env?: Record<string, string>;
    models?: string[];
  };
}

export interface ProvidersConfig {
  version: number;
  providers: Provider[];
}
