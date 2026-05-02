# Feature: Parallel Launch（Profile 模式）

## 背景

proteus 现有的 `switch` 命令通过覆写 `~/.claude/settings.json` 来切换 provider，同一时间只能激活一个 provider。

这个机制适合“长期默认切换”，但不适合“同一时间并行跑多个 Claude配置”的场景：

- 用 DeepSeek 跑日常编码，用 Anthropic 官方跑需要高质量输出的任务
- 同一个项目开两个 Claude，一个写代码，一个做 review
- 对比不同 provider 在同一问题上的输出

## 目标

新增 `launch` 子命令，按 **profile** 启动 Claude：

- 从 profile 解析 provider
- 将 provider 的配置用于构造 `execEnv`
- `syscall.Exec` 替换进程启动 `claude`

从而实现多个 Claude 进程并行运行、互不干扰。

同时，**保留现有 `switch` 的写入逻辑不变**：`switch` 继续覆写 `~/.claude/settings.json`，用于全局持久切换。

## 命令职责

- `proteus switch`：写入 `~/.claude/settings.json`（全局、持久）
- `proteus launch`：构造进程 `execEnv` 并 exec（当前进程、临时）

两个命令并存，职责明确，互不替代。

## 核心机制

`launch` 不使用 `CLAUDE_CONFIG_DIR` 隔离，不写任何配置文件；仅在当前进程内构造目标环境变量，然后 `syscall.Exec` 启动 `claude`：

1. 读取 profile 配置
2. 解析其绑定的 provider，并读取 `provider.claude.env`
3. 基于当前进程环境与 provider 的 `claude.env` 生成最终 `execEnv`（value 经过 `os.ExpandEnv()`）
4. 组装为 `[]string{"K=V"}` 传给 `syscall.Exec`
5. `syscall.Exec` 替换当前进程

Claude Code 启动时会读取进程环境变量中的 `ANTHROPIC_API_KEY`、`ANTHROPIC_BASE_URL`、`ANTHROPIC_MODEL` 等，优先级高于 `settings.json` 的 `env` 字段。

## 用法

```bash
# 按 profile 启动 Claude
proteus launch coding-fast

# 预览会设置哪些环境变量，不实际启动
proteus launch coding-fast --dry-run

# 列出支持 launch 的 profile
proteus launch --list
```

## 配置

在 `providers.yaml` 中新增 `profiles`，profile 仅负责命名和绑定 provider，运行时配置统一来自 provider。

```yaml
version: 1

providers:
  deepseek:
    name: DeepSeek
    claude:
      env:
        ANTHROPIC_API_KEY: "${DEEPSEEK_API_KEY}"
        ANTHROPIC_BASE_URL: "https://api.deepseek.com/v1"
      models:
        - deepseek-chat
        - deepseek-reasoner

  anthropic:
    name: Anthropic
    claude:
      env:
        ANTHROPIC_API_KEY: "${ANTHROPIC_API_KEY}"
      models:
        - claude-opus-4-6
        - claude-sonnet-4-6

profiles:
  coding-fast:
    provider: deepseek

  review-strict:
    provider: anthropic
```

env value 支持 `${VAR}` 语法引用宿主环境变量，API key 不需要明文写入配置文件。

## env 合并规则

`launch` 最终生效 env 优先级（后者覆盖前者）：

1. 当前进程 `os.Environ()`
2. `provider.claude.env`

## 与 switch 的区别

| | `proteus switch` | `proteus launch` |
|---|---|---|
| 机制 | 覆写 `~/.claude/settings.json` | 构造 `execEnv` 并 exec `claude` |
| 并行 | 不支持 | 支持，每个终端独立 |
| 持久化 | 全局生效，重启后保留 | 仅当前进程，退出即消失 |
| 选择粒度 | provider | profile（间接绑定 provider） |
| 适用场景 | 长期切换默认 provider | 临时并行启动多个 Claude |

## 实现要点

### 目录结构

```text
internal/
  app/
    launch.go       # Launch 入口，参数解析与 --dry-run
  launcher/
    exec.go         # profile 解析、env 合并、execEnv 构造、syscall.Exec
```

### env 注入

推荐实现：先基于 `os.Environ()` 构造目标 env map，叠加 `provider.claude.env`（value 做 `os.ExpandEnv()`），再组装为 `[]string{"K=V"}` 传给 `syscall.Exec`。

```go
baseEnv := environToMap(os.Environ())
for k, v := range providerEnv {
    baseEnv[k] = os.ExpandEnv(v)
}
execEnv := mapToEnviron(baseEnv)
syscall.Exec(claudePath, []string{"claude"}, execEnv)
```

说明：

- 不依赖 `os.Setenv()` 改写当前进程环境，减少副作用窗口
- 语义更直接：最终传给 `Exec` 的 env 即最终生效 env
- 行为与文档中的优先级规则一致（`provider.claude.env` 覆盖当前进程）

`launch` 不写任何文件，不创建任何目录，不修改 `~/.claude/settings.json`。

### --dry-run 输出

默认脱敏展示敏感变量（如 `ANTHROPIC_API_KEY`），避免在终端历史或录屏中泄漏密钥。

```text
Profile:  coding-fast
Provider: deepseek (DeepSeek)
Command:  /usr/local/bin/claude

Env:
  ANTHROPIC_API_KEY  = sk-***...***
  ANTHROPIC_BASE_URL = https://api.deepseek.com/v1
  ANTHROPIC_MODEL    = deepseek-chat
```

可选增强：增加 `--show-secrets`（默认关闭），仅在用户明确需要排查时显示原值。

### 错误处理

- profile 不存在：提示可用 profile 列表
- profile 引用的 provider 不存在：明确报错
- `claude` 不在 PATH：明确报错，提示安装（实现建议使用 `exec.LookPath("claude")`）
- provider 与 profile 合并后 env 为空：报错，提示检查配置
- `os.ExpandEnv` 展开后值为空（引用的环境变量未设置）：warn 提示，不阻断启动
- 关键变量（如 `ANTHROPIC_API_KEY`）为空：输出更醒目的 `WARN[critical]`，并提示可能导致鉴权失败（仍可继续启动）

### 可观测性

- `launch --list` 建议显示：profile、provider、provider 是否存在、关键 env 是否可展开
- `--dry-run` 建议显示 `claude` 解析后的绝对路径，便于 PATH 排查

## 不在范围内

- 同时启动多个 profile（`proteus launch coding-fast review-strict`）：第一版不做，`syscall.Exec` 语义不适合批量启动
- Windows 支持：`syscall.Exec` 是 Unix-only，Windows 需要另一套实现，暂不考虑
- 自动打开新终端窗口：launch 在当前终端 exec，不调用 Terminal.app / osascript
