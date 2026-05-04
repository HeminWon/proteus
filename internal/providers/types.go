package providers

type Provider struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Claude struct {
		Env    map[string]string `yaml:"env"`
		Models []string          `yaml:"models,omitempty"`
	} `yaml:"claude"`
	Codex *struct {
		Env    map[string]string `yaml:"env,omitempty"`
		Models []string          `yaml:"models,omitempty"`
	} `yaml:"codex,omitempty"`
}

type Profile struct {
	Provider string   `yaml:"provider"`
	Runner   string   `yaml:"runner,omitempty"`
	Args     []string `yaml:"args,omitempty"`
}

type ProvidersConfig struct {
	Version   int                `yaml:"version"`
	Providers []Provider         `yaml:"providers"`
	Profiles  map[string]Profile `yaml:"profiles,omitempty"`
}

type LoadProvidersResult struct {
	Config    ProvidersConfig
	ConfigDir string
}
