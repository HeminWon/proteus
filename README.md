# Proteus

Proteus is a CLI for switching Claude Code provider settings and launching profile-isolated sessions.

## Install

### Homebrew

```bash
brew tap HeminWon/proteus https://github.com/HeminWon/proteus
brew install proteus
```

### From source

```bash
go run ./cmd/proteus --help
```

## Configuration

1. Copy the example config:

```bash
cp configs/providers.example.yaml ~/.config/proteus/providers.yaml
```

2. Edit `~/.config/proteus/providers.yaml` and set your provider env values.

3. (Optional) Use a custom config directory in `~/.config/proteus/config.json`:

```json
{
  "config_dir": "~/my-providers"
}
```

## Commands

```bash
proteus --help
proteus list
proteus validate [--provider <id>] [--concurrency <n>]
proteus switch <provider-id|provider-name> [--dry-run]
proteus launch <profile> [--dry-run]
proteus launch --list
```

## Launch behavior

- `switch` writes global settings to `~/.claude/settings.json`.
- `launch` does **not** write global settings.
- `launch` writes profile-private settings and starts the configured runner with profile env.
- `profile.runner` must be an executable name only (for example `claude` or `codex`).
- Put runner flags in `profile.args` instead of `runner`.

## Notes

- Do not commit real tokens in `providers.yaml`.
- If both `ANTHROPIC_AUTH_TOKEN` and `ANTHROPIC_API_KEY` are empty, authentication will fail.
- Use `proteus <command> --help` for command-specific usage.
