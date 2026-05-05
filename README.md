<p align="center">
  <img src="assets/logo.png" alt="Proteus logo" width="180" />
</p>

<h1 align="center">Proteus</h1>

<p align="center">
  A Go CLI to switch Claude Code providers and launch profile-isolated sessions.
</p>

<p align="center">
  <a href="README.zh-CN.md">中文文档</a>
</p>

## Background

`switch` is useful when you only need one global provider/model at a time.

When you need parallel workflows (for example coding + review + experiments), Proteus `launch` lets you run multiple Claude Code terminals simultaneously with different profiles, providers, and model presets, without overwriting each other's runtime state.

## Typical Scenarios

- **Single global context**: You only need one provider/model at a time. Use `switch` to update global Claude settings quickly.
- **Parallel workflows**: You need to use models from different providers (for example DeepSeek and GLM) while running tasks in parallel terminals.
  - Terminal A: `proteus launch deepseek`
  - Terminal B: `proteus launch glm`
  - Terminal C: `proteus launch anthropic`

## `switch` vs `launch`

| Command | Writes global `~/.claude/settings.json` | Profile/session isolation | Best for |
| --- | --- | --- | --- |
| `proteus switch <provider>` | Yes | No | Fast global provider/model switch |
| `proteus launch <profile>` | No | Yes | Running multiple isolated Claude Code sessions in parallel |

## Features

- Manage multiple Claude-compatible providers in one config file.
- Switch active provider by writing `~/.claude/settings.json`.
- Launch profile-isolated sessions without mutating global settings.
- Run multiple Claude Code terminals in parallel, each with different profiles/providers/model presets.
- Sync shared Claude config entries (`commands`, `skills`, `plugins`, `agents`, `ide`) into profile config dir.
- Validate provider configuration with live HTTP checks.

## Requirements

- Go `1.22+` (for building/running from source)
- Claude Code installed if you use `launch` with `runner: claude`

## Installation

### Homebrew

```bash
brew tap HeminWon/proteus https://github.com/HeminWon/proteus
brew install proteus
```

### Build from source

```bash
go build -o dist/proteus ./cmd/proteus
./dist/proteus --help
```

## Quick Start

1. Create your provider config:

```bash
cp configs/providers.example.yaml ~/.config/proteus/providers.yaml
```

2. Edit `~/.config/proteus/providers.yaml` and fill your token/env values.

3. Validate configuration:

```bash
proteus validate
```

4. Switch provider globally:

```bash
proteus switch anthropic
```

5. Launch an isolated profile session:

```bash
proteus launch default
```

## Configuration

### Providers file

Default path:

```text
~/.config/proteus/providers.yaml
```

Minimal example:

```yaml
version: 1
providers:
  - id: anthropic
    name: Anthropic Official
    claude:
      env:
        ANTHROPIC_AUTH_TOKEN: "change-me"
        ANTHROPIC_BASE_URL: "https://api.anthropic.com"

profiles:
  default:
    provider: anthropic
    runner: claude
    args: []
    share_claude_md: false
```

### Optional config directory override

File:

```text
~/.config/proteus/config.json
```

Content:

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

## Behavior Notes

- `switch` persists provider env into global `~/.claude/settings.json`.
- `launch` does not write global settings; it writes profile-private settings and starts `profile.runner` with `profile.args`.
- `profile.runner` must be an executable name (for example `claude` or `codex`).
- `share_claude_md` is `false` by default.

## Security Notes

- Never commit real tokens in `providers.yaml`.
- If both `ANTHROPIC_AUTH_TOKEN` and `ANTHROPIC_API_KEY` are empty, authentication will fail.

## Development

```bash
go test ./...
```

Or use `just` tasks:

```bash
just build
just run
just list
just validate
```

## Contributing

Contributions are welcome. Please open an issue or pull request with:

- a clear problem statement,
- reproducible steps (for bug reports),
- and expected behavior.

## License

This project is licensed under the [MIT License](LICENSE).
