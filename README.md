# Proteus

Proteus 用于快速切换 Claude Code 的 provider 配置（如 Anthropic / OpenRouter / AIClient2API）。

## 安装

### Homebrew

```bash
brew tap HeminWon/proteus https://github.com/HeminWon/proteus
brew install proteus
```

### 开发模式

```bash
go run ./cmd/proteus --help
```

## 配置

1. 复制示例配置：

```bash
cp configs/providers.example.yaml ~/.config/proteus/providers.yaml
```

2. 编辑 `~/.config/proteus/providers.yaml`，填入你的 provider 信息（token、base URL、models 等）。

3. （可选）自定义配置目录：编辑 `~/.config/proteus/config.json`

```json
{
  "config_dir": "~/my-providers"
}
```

## 常用命令

```bash
proteus --help                              # 显示帮助
proteus list                                # 列出 provider
proteus switch <provider-name>              # 切换 provider
proteus switch <provider-name> --dry-run    # 预览切换
proteus validate                            # 校验配置与连通性
```

## 注意

- `providers.yaml` 请勿提交包含真实 token 的版本。
- 切换时会更新 `~/.claude/settings.json` 的 `env` 和 `availableModels`。
