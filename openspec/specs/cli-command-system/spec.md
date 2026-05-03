# cli-command-system Specification

## Purpose
TBD - created by archiving change cli-optimization. Update Purpose after archive.
## Requirements
### Requirement: Unified subcommand-oriented command system
The CLI MUST use subcommand-oriented invocation for primary actions (`list`, `validate`, `switch`, `launch`).

#### Scenario: Use canonical list command
- **WHEN** user runs `proteus list`
- **THEN** the system MUST execute provider listing behavior

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

