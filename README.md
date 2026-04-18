# Proteus

> *"他能变成任何形态——狮子、蛇、豹、猪、流水、火焰、大树。但只要你抓住他不放，他终将现出本来面目，告诉你真相。"*
> — 荷马《奥德赛》

---

## 起源

在希腊神话中，Proteus 是海洋之神，波塞冬的牧人，掌管着爱琴海的海豹群。他拥有两种能力：**预言**与**变形**。

他不轻易示人。想从他口中得到答案，你必须在正午趁他熟睡时抓住他，任凭他变幻千种形态，绝不松手——直到他精疲力竭，恢复原形，才会开口说出真相。

这个工具借用了他的名字，因为它做的事情与 Proteus 的本质如出一辙：

**在不同形态之间切换，但内核始终如一。**

Claude Code 是那个"内核"。Anthropic、OpenRouter、AIClient2API——这些 provider 是它的"形态"。Proteus 让你在这些形态之间自由穿梭，而不必每次都去手动翻找配置文件。

---

## 它解决什么问题

AI CLI 工具的 provider 切换，本质上是修改几个环境变量。但这件事足够频繁、足够烦人：

- 官方 API 额度用完，临时切到中转
- 测试不同 provider 的模型表现
- 在不同网络环境下使用不同入口

每次都要打开 `~/.claude/settings.json`，找到 `env` 字段，手动替换 token 和 base URL，再保存。重复、易错、无聊。

Proteus 把这件事变成一行命令。

---

## 安装

### Homebrew（Formula）

> `proteus` 是 CLI 工具，使用 Homebrew Formula（不是 cask）。

安装方式：

```bash
brew tap HeminWon/proteus https://github.com/HeminWon/proteus
brew install proteus
```

升级：

```bash
brew update
brew upgrade proteus
```

说明：
- 本仓库自带 tap 与 `Formula/proteus.rb`
- 支持平台：macOS (arm64)、Linux (arm64/x64)
- 发布 tag（`v*`）后会自动更新 Formula 中的 version/url/sha256
- 若自动更新失败，可在 Actions 手动触发 `update-homebrew-formula`（`workflow_dispatch`）补跑
- 发布二进制使用 Bun 稳定版（`latest`）

### 开发模式

```bash
npm ci
npx tsx src/cli/index.ts --list
```

也可以直接用 package script：

```bash
npm run list
npm run validate
```

### 配置示例

仓库内提供了 `configs/providers.example.yaml`，复制后填入你自己的 provider：

```bash
cp configs/providers.example.yaml ~/.config/proteus/providers.yaml
```

### 打包为二进制

```bash
npm run build:bun
```

生成的 `proteus` 二进制可直接运行，或复制到 PATH：

```bash
cp dist/proteus /usr/local/bin/proteus
proteus --list
```

注意：在较新的 macOS 上，Bun `--compile` 生成的单文件二进制可能被系统安全策略拦截。开发和自用场景下，优先直接运行 `tsx` 版本，或用一个 shell wrapper 调 `npx tsx src/cli/index.ts`。

已知兼容性提示：Bun `1.3.12` 在部分 macOS 环境可能导致编译产物运行时被系统直接 `killed`。建议使用 `1.3.13+` 后再执行 `build:bun`。

---

## 当前能力

```bash
proteus                             # 列出所有 provider
proteus anthropic                   # 切换到指定 provider
proteus --validate                  # 校验 providers.yaml + curl 验证 token/base URL
proteus anthropic --dry-run         # 预览切换变更，不写入
proteus --help                      # 查看完整帮助
```

`--list`（或不带参数）的结果会直接标记当前激活 provider（`◀ active`）。
`--validate` 会输出每个 provider 的 live 校验结果与延迟（`latency=xxms`）。

参数规则（避免误用）：

- `--list` / `--validate` / `--help` 互斥，不能组合
- `--dry-run` 只能用于切换 provider（如 `npx tsx src/cli/index.ts openrouter --dry-run`）
- 未知参数会直接报错，并提示使用 `--help`

配置统一写在 `providers.yaml`，切换时只更新 `~/.claude/settings.json` 中的 `env` 和 `availableModels`。
不要把包含真实 token 的 `providers.yaml` 提交到仓库。

### 配置路径（XDG 规范）

- `~/.config/proteus/config.json` — 主配置，可指定 `config_dir` 字段自定义 providers.yaml 路径
- `~/.config/proteus/providers.yaml` — providers 配置（默认位置）
- `~/.cache/proteus/cache.json` — 当前激活 provider 缓存

`config.json` 示例：

```json
{
  "config_dir": "~/my-providers"
}
```

未配置时默认使用 `~/.config/proteus/`。

当前激活 provider 不再写回 `providers.yaml`，而是缓存到 `~/.cache/proteus/cache.json`。

首次没有 cache 时，active 状态为未设置（unset），不会默认选中第一个 provider。

切换流程包含以下保护措施：

- 默认 `overwrite-env`：切换时用 provider 的 `claude.env` 全量覆盖 live `env`
- `--dry-run`：先展示变更计划（新增/更新/删除字段）
- 切换前自动备份 `settings.json` 到 `~/.claude/proteus-backups/`
- 自动轮换备份，仅保留最近 10 份
- 当前激活状态通过本地 cache 记录，不修改 `providers.yaml`
- 原子写入，写入失败时尝试用备份恢复

---

## 后续规划

**近期**

- 支持 `proteus add` 交互式添加新 provider

**中期**

- 扩展到其他 AI CLI 工具：Codex、Gemini CLI、OpenCode
- 支持"环境组"概念：一条命令同时切换多个工具的 provider
- 切换历史记录，支持 `proteus back` 回退到上一个 provider

**长期**

- 脱离 `npx tsx`，编译为单一可执行文件，全局安装
- 支持 profile 概念：工作 / 个人 / 测试 等场景预设
- 与 heminSpec 的 link-manifest 体系集成，作为环境管理的一部分

---

*Proteus 不是一个大工具。它只做一件事，并且把它做好。*
