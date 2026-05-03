package core

type CliAction string

const (
	ActionList     CliAction = "list"
	ActionValidate CliAction = "validate"
	ActionSwitch   CliAction = "switch"
	ActionLaunch   CliAction = "launch"
	ActionHelp     CliAction = "help"
)

type CliOptions struct {
	Action              CliAction
	ProviderInput       string
	ProfileInput        string
	DryRun              bool
	ListLaunch          bool
	HelpCommand         string
	ValidateProvider    string
	ValidateConcurrency int
}
