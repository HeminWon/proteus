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

## 功能特性

- 在一个配置文件中管理多个 Claude 兼容 Provider。
- 通过写入 `~/.claude/settings.json` 切换当前全局 Provider。
- 启动 Profile 隔离会话，不污染全局设置。
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

1. 创建 Provider 配置文件：

```bash
cp configs/providers.example.yaml ~/.config/proteus/providers.yaml
```

2. 编辑 `~/.config/proteus/providers.yaml`，填入你的 token/env。

3. 校验配置：

```bash
proteus validate
```

4. 全局切换 Provider：

```bash
proteus switch anthropic
```

5. 启动隔离 Profile 会话：

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
    args: []
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

## 行为说明

- `switch` 会将 Provider env 持久化到全局 `~/.claude/settings.json`。
- `launch` 不写全局设置；会写入 Profile 私有 settings，并用 `profile.runner` + `profile.args` 启动进程。
- `profile.runner` 必须是可执行文件名（例如 `claude`、`codex`）。
- `share_claude_md` 默认值为 `false`。

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

当前仓库尚未包含 `LICENSE` 文件。若计划面向开源社区分发，建议补充许可证（如 MIT/Apache-2.0）。
