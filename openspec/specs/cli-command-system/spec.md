# cli-command-system Specification

## Purpose
TBD - created by archiving change cli-optimization. Update Purpose after archive.
## Requirements
### Requirement: Unified subcommand-oriented command system
The CLI MUST use subcommand-oriented invocation for primary actions (`list`, `validate`, `switch`, `launch`, `doctor`) and MUST support transitional compatibility routing for legacy global flags.

#### Scenario: Use canonical list command
- **WHEN** user runs `proteus list`
- **THEN** the system MUST execute provider listing behavior
- **AND** output MUST be equivalent to legacy list behavior

#### Scenario: Route legacy global flag to canonical command
- **WHEN** user runs legacy `proteus --list`
- **THEN** the system MUST execute the same behavior as `proteus list`
- **AND** the system MUST print a deprecation warning indicating the canonical command

### Requirement: Canonical validate command with scoped options
The CLI MUST provide `proteus validate` as the canonical validation entrypoint and MUST allow provider-scoped and concurrency-scoped validation options.

#### Scenario: Validate all providers
- **WHEN** user runs `proteus validate`
- **THEN** the system MUST validate configured providers

#### Scenario: Validate a single provider
- **WHEN** user runs `proteus validate --provider anthropic`
- **THEN** the system MUST validate only provider `anthropic`

#### Scenario: Validate with explicit concurrency
- **WHEN** user runs `proteus validate --concurrency 10`
- **THEN** the system MUST apply concurrency level `10` for validation operations

