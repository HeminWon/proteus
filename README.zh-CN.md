<p align="center">
  <img src="assets/logo.png" alt="Proteus logo" width="180" />
</p>

<h1 align="center">Proteus</h1>

<p align="center">
  一个用于切换 Claude Code Provider 并启动 Profile 隔离会话的 Go CLI。
</p>

<p align="center">
  <a href="README.md">English Documentation</a>
</p>

## `switch` 与 `launch` 对比

### `proteus switch <provider>`
- 全局设置：会写入 `~/.claude/settings.json`
- 隔离性：无（全局生效）
- 适用场景：快速切换当前默认 Provider/模型

### `proteus launch <profile>`
- 全局设置：不写入全局 settings
- 隔离性：有（Profile/会话隔离）
- 运行方式：使用 `profile.runner` + `profile.args` 启动（`profile.runner` 必须是可执行名，如 `claude`、`codex`）
- 默认值：`share_claude_md` 为 `false`
- 适用场景：并行运行多个不同 Provider 的会话（如 DeepSeek / GLM / Anthropic）

## 功能特性

- 在一个配置文件中管理多个 Claude 兼容 Provider。
- 同时支持全局切换（`switch`）与 Profile 隔离并行会话（`launch`）。
- 将共享 Claude 配置项（`commands`、`skills`、`plugins`、`agents`、`ide`）同步到 Profile 配置目录。
- 支持带实时 HTTP 检查的配置校验。

## 环境要求

- Go `1.22+`（源码构建/运行）
- 如果使用 `launch` 且 `runner: claude`，需要本机已安装 Claude Code

## 安装

### Homebrew

```bash
brew tap HeminWon/proteus https://github.com/HeminWon/proteus
brew install proteus
```

### 源码构建

```bash
go build -o dist/proteus ./cmd/proteus
./dist/proteus --help
```

## 快速开始

1. 创建配置：

```bash
cp configs/providers.example.yaml ~/.config/proteus/providers.yaml
```

2. 在 `~/.config/proteus/providers.yaml` 填入 token/env。

3. 校验：

```bash
proteus validate
```

4. 全局切换：

```bash
proteus switch anthropic
```

5. 隔离启动：

```bash
proteus launch default
```

## 配置说明

### Providers 文件

默认路径：

```text
~/.config/proteus/providers.yaml
```

最小示例：

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
    args:
      - --dangerously-skip-permissions
    share_claude_md: false
```

### 可选配置目录覆盖

文件：

```text
~/.config/proteus/config.json
```

内容：

```json
{
  "config_dir": "~/my-providers"
}
```

## 命令列表

```bash
proteus --help
proteus list
proteus validate [--provider <id>] [--concurrency <n>]
proteus switch <provider-id|provider-name> [--dry-run]
proteus launch <profile> [--dry-run]
proteus launch --list
```

## 安全提示

- 不要在 `providers.yaml` 中提交真实 token。
- 若 `ANTHROPIC_AUTH_TOKEN` 与 `ANTHROPIC_API_KEY` 均为空，会导致认证失败。

## 开发

```bash
go test ./...
```

或使用 `just`：

```bash
just build
just run
just list
just validate
```

## 贡献

欢迎提交 Issue 或 PR，建议包含：

- 明确的问题描述，
- 可复现步骤（Bug 场景），
- 期望行为。

## 许可证

本项目基于 [MIT License](LICENSE) 开源。
