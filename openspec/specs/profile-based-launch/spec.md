## Purpose
Enable launching Claude by profile with safe runtime env resolution, clear diagnostics, and no persistent global settings mutation.

## Requirements

### Requirement: Launch Claude by profile without persisting settings
The system MUST provide a `proteus launch <profile>` command that resolves the target provider from the profile and starts `claude` via process replacement (`syscall.Exec`) using runtime environment injection, without writing any Claude settings file.

#### Scenario: Launch with a valid profile
- **WHEN** user runs `proteus launch coding-fast` and profile `coding-fast` maps to an existing provider
- **THEN** system resolves provider config, builds final exec environment, and replaces current process with `claude`
- **AND** system does not create or modify `~/.claude/settings.json`

#### Scenario: Profile does not exist
- **WHEN** user runs `proteus launch unknown-profile`
- **THEN** system MUST return a clear error indicating the profile is missing
- **AND** output available profiles to help correction

#### Scenario: Referenced provider does not exist
- **WHEN** profile exists but its `provider` field references a non-existent provider
- **THEN** system MUST return a clear error indicating invalid profile-provider binding

### Requirement: Merge runtime environment with provider precedence
The system MUST construct launch environment by starting from current process environment and then overriding with `provider.claude.env` values after `os.ExpandEnv` expansion.

#### Scenario: Provider env overrides existing process env
- **WHEN** current environment contains key `ANTHROPIC_BASE_URL=A` and provider env sets `ANTHROPIC_BASE_URL=B`
- **THEN** final exec environment MUST use `ANTHROPIC_BASE_URL=B`

#### Scenario: Provider env references host environment variable
- **WHEN** provider env value is `${DEEPSEEK_API_KEY}` and host variable `DEEPSEEK_API_KEY` is set
- **THEN** final exec environment MUST contain expanded value

#### Scenario: Expansion result is empty
- **WHEN** provider env references an unset host variable and expansion result is empty
- **THEN** system MUST emit a warning
- **AND** launch MAY continue

### Requirement: Provide dry-run observability with secret masking
The system MUST support `proteus launch <profile> --dry-run` to show resolved launch details without executing `claude`, and MUST mask sensitive values by default.

#### Scenario: Dry-run displays launch context
- **WHEN** user runs `proteus launch coding-fast --dry-run`
- **THEN** output MUST include profile name, resolved provider, resolved claude executable path, and effective key env entries
- **AND** command execution MUST stop before `syscall.Exec`

#### Scenario: Sensitive env values are masked by default
- **WHEN** dry-run output includes `ANTHROPIC_API_KEY`
- **THEN** value MUST be masked to prevent full secret exposure

### Requirement: List launch-capable profiles
The system MUST support `proteus launch --list` to enumerate configured profiles and their provider mapping for launch usage.

#### Scenario: List profiles
- **WHEN** user runs `proteus launch --list`
- **THEN** system MUST output all profiles
- **AND** each entry MUST include profile name and referenced provider

### Requirement: Validate launch prerequisites and fail clearly
The system MUST validate critical launch prerequisites and provide explicit errors or warnings.

#### Scenario: Claude executable missing from PATH
- **WHEN** `exec.LookPath("claude")` fails
- **THEN** system MUST return a clear error indicating `claude` is not available in PATH

#### Scenario: Effective provider env is empty
- **WHEN** merged launch env has no provider-contributed keys after resolution
- **THEN** system MUST return an error indicating invalid/empty launch configuration

#### Scenario: Critical auth variable is empty
- **WHEN** effective `ANTHROPIC_API_KEY` is empty
- **THEN** system MUST emit a prominent critical warning
- **AND** launch MAY continue
