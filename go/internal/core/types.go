package core

type CacheData struct {
	Active *struct {
		Claude string `json:"claude,omitempty"`
	} `json:"active,omitempty"`
}

type CliAction string

const (
	ActionList     CliAction = "list"
	ActionValidate CliAction = "validate"
	ActionSwitch   CliAction = "switch"
	ActionHelp     CliAction = "help"
)

type CliOptions struct {
	Action        CliAction
	ProviderInput string
	DryRun        bool
}
