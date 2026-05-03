# cli-help-and-suggestions Specification

## Purpose
TBD - created by archiving change cli-optimization. Update Purpose after archive.
## Requirements
### Requirement: Subcommand-scoped help output
The CLI MUST provide context-specific help output for supported subcommands and MUST preserve global help for top-level usage.

#### Scenario: Show global help
- **WHEN** user runs `proteus --help`
- **THEN** the system MUST show top-level usage and available commands
- **AND** the output MUST focus on global command discovery

#### Scenario: Show switch help
- **WHEN** user runs `proteus switch --help`
- **THEN** the system MUST show usage, flags, and examples specific to `switch`

#### Scenario: Show launch help
- **WHEN** user runs `proteus launch --help`
- **THEN** the system MUST show usage, flags, and examples specific to `launch`
- **AND** the output MUST clarify that `launch` does not persist settings files

### Requirement: Typo suggestion for recoverable CLI input errors
The CLI MUST provide suggestion text for recoverable input errors by matching user input against known candidates for commands, providers, profiles, and flags.

#### Scenario: Suggest command name for unknown command
- **WHEN** user runs an unknown command `swith`
- **THEN** the system MUST return an error message
- **AND** the message MUST include suggestion `switch` when similarity threshold is met

#### Scenario: Suggest provider identifier for switch target
- **WHEN** user runs `proteus switch anthopic` and provider `anthropic` exists
- **THEN** the system MUST return an error message
- **AND** the message MUST include suggestion `anthropic` when similarity threshold is met

#### Scenario: Suggest profile name for launch target
- **WHEN** user runs `proteus launch defualt` and profile `default` exists
- **THEN** the system MUST return an error message
- **AND** the message MUST include suggestion `default` when similarity threshold is met

#### Scenario: Suggest known flag for unknown option
- **WHEN** user runs `proteus switch --dryun`
- **THEN** the system MUST return an unknown option error
- **AND** the message MUST include suggestion `--dry-run` when similarity threshold is met

