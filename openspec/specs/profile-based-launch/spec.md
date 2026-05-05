## MODIFIED Requirements

### Requirement: Launch Claude by profile without persisting settings
The system MUST provide a `proteus launch <profile>` command that resolves the target provider from the profile, prepares a profile-specific `CLAUDE_CONFIG_DIR`, and starts `claude` via process replacement (`syscall.Exec`) using runtime environment injection. The launch flow MUST only write files inside the profile-specific config directory and MUST NOT create or modify `~/.claude/settings.json`.

#### Scenario: Launch with a valid profile
- **WHEN** user runs `proteus launch coding-fast` and profile `coding-fast` maps to an existing provider
- **THEN** system resolves provider config, prepares the profile config directory, and replaces current process with `claude`
- **AND** system writes settings only to the profile-private settings path under that directory
- **AND** system does not create or modify `~/.claude/settings.json`

#### Scenario: Profile does not exist
- **WHEN** user runs `proteus launch unknown-profile`
- **THEN** system MUST return a clear error indicating the profile is missing
- **AND** output available profiles to help correction

#### Scenario: Referenced provider does not exist
- **WHEN** profile exists but its `provider` field references a non-existent provider
- **THEN** system MUST return a clear error indicating invalid profile-provider binding

### Requirement: Synchronize shareable Claude config entries into profile config directory
The system MUST prepare the profile config directory before launch by synchronizing a fixed whitelist of shareable entries from `~/.claude` into the profile-specific `CLAUDE_CONFIG_DIR` using symlinks. The whitelist MUST include `commands/`, `skills/`, `plugins/`, `agents/`, and `ide/`. Runtime isolation MUST be achieved by only synchronizing this whitelist, not by scanning all global entries and filtering a blacklist.

#### Scenario: Shareable source entry exists
- **WHEN** a whitelisted entry exists under `~/.claude`
- **THEN** system MUST create a symlink at the matching path under the profile config directory pointing to the global entry

#### Scenario: Shareable source entry missing
- **WHEN** a whitelisted entry does not exist under `~/.claude`
- **THEN** system MUST skip that entry without failing launch

#### Scenario: Existing symlink already matches target
- **WHEN** the destination path already exists as a symlink to the expected global entry
- **THEN** system MUST leave it unchanged and continue launch

#### Scenario: Destination path is a non-matching file or directory
- **WHEN** the destination path already exists but is not the expected symlink target
- **THEN** system MUST fail launch with a clear error identifying the conflicting path

### Requirement: Keep runtime state isolated across profiles
The system MUST NOT symlink runtime state, cache, history, or session data from `~/.claude` into the profile-specific `CLAUDE_CONFIG_DIR`.

#### Scenario: Runtime state exists in global Claude directory
- **WHEN** entries such as `history.jsonl`, `sessions/`, `session-env/`, `transcripts/`, `tasks/`, `todos/`, `plans/`, `projects/`, `telemetry/`, `statsig/`, `stats-cache.json`, `cache/`, `paste-cache/`, `file-history/`, `shell-snapshots/`, `debug/`, `backups/`, or `proteus-backups/` exist under `~/.claude`
- **THEN** system MUST NOT create symlinks for those entries in the profile config directory

#### Scenario: Launch creates profile-private runtime state later
- **WHEN** Claude writes runtime files after launch
- **THEN** those files MUST remain inside the profile-specific `CLAUDE_CONFIG_DIR` rather than referencing global runtime state

### Requirement: Support profile-level CLAUDE.md boolean sharing policy
The system MUST allow each profile to control through a boolean configuration field whether `CLAUDE.md` is shared from `~/.claude/CLAUDE.md` into the profile config directory. The default behavior MUST be not to share `CLAUDE.md`.

#### Scenario: Profile enables shared CLAUDE.md
- **WHEN** profile configuration indicates that `CLAUDE.md` should be shared and the global file exists
- **THEN** system MUST create or preserve a symlink from the profile config directory to `~/.claude/CLAUDE.md`

#### Scenario: Profile disables shared CLAUDE.md
- **WHEN** profile configuration indicates that `CLAUDE.md` should not be shared
- **THEN** system MUST NOT create a symlink for `CLAUDE.md`
- **AND** launch MUST leave any profile-local `CLAUDE.md` untouched
- **AND** launch MUST NOT delete or rewrite an existing `CLAUDE.md` path automatically

#### Scenario: Profile disables shared CLAUDE.md after a shared symlink already exists
- **WHEN** profile configuration indicates that `CLAUDE.md` should not be shared
- **AND** the destination path already exists as a symlink to `~/.claude/CLAUDE.md`
- **THEN** system MUST leave that path unchanged
- **AND** system SHOULD report that state in diagnostics as an existing path retained under disabled sharing

### Requirement: Expose sync results in dry-run output
The system MUST include profile config synchronization results in `proteus launch <profile> --dry-run` output.

#### Scenario: Dry-run shows pending sync actions
- **WHEN** user runs `proteus launch coding-fast --dry-run`
- **THEN** output MUST include the profile config directory and each shareable entry's sync status such as `linked`, `reused`, `skipped-missing`, `conflict`, `disabled`, or `disabled-existing`
- **AND** command execution MUST stop before `syscall.Exec`

#### Scenario: Dry-run has no filesystem side effects
- **WHEN** user runs `proteus launch coding-fast --dry-run`
- **THEN** system MUST NOT create the profile directory
- **AND** system MUST NOT create, delete, or modify any symlink or file
- **AND** system MUST NOT write profile-private `settings.json`
