# QA Helper Skills

`qa-hepler-skills` 是一个同时面向 Codex plugin 和 Claude plugin 的 QA 工作流插件。

## 能力

- 需求分析与需求自动补全，支持来源选项：`1. 手动输入`、`2. Jira`、`3. TM`
- 通过 `--jiraid <JIRA-ID>` 从 Jira MCP 获取需求详情后进行分析
- 功能测试用例设计
- 测试用例导出为 XMind
- 测试用例导出为 Excel
- 用例导入导出支持选项：`1. 通过 XMind 输出`、`2. 通过 Excel 输出`
- 通过 `cml-mcp` 查询、新增、查看、更新和导入 MeterSphere/CML 功能用例
- 端到端执行 `需求分析 -> 功能用例设计 -> 用例导入`
- 端到端执行时，每一步 skill 结果写入 `tmp/.cache` 并等待用户确认后继续

## Codex Plugin

Codex 插件入口：

```text
.codex-plugin/plugin.json
```

插件会加载：

- `skills/`
- `.mcp.json`

## Claude Plugin

Claude 插件入口：

```text
.claude-plugin/plugin.json
```

可用命令：

- `/requirement`
- `/testcase`
- `/case-export`
- `/case-import`
- `/requirement-case`

## MCP 配置

共享 MCP 配置位于：

```text
.mcp.json
```

启动 `cml-mcp` 前需要配置 MeterSphere/CML 环境变量：

```bash
export CML_BASE_URL="http://localhost:8081"
export CML_ACCESS_KEY="your-access-key"
export CML_SIGNATURE="your-signature"
export CML_PROJECT_ID="project-id"
```

也可以在项目根目录或 `mcp/cml-mcp` 目录放置 `.env` 文件，`cml-mcp` 启动时会自动读取；已存在的系统环境变量和命令行参数优先级更高。

Cookie 认证可额外配置：

```bash
export CML_COOKIE="SESSION=..."
```

日志级别：

```bash
export CML_LOG_LEVEL="info"
```

`cml-mcp` 日志输出到 `stderr`，不会污染 MCP `stdout` 协议流。

## Jira 需求分析

需求分析 skill 触发时支持需求描述来源选择：

```text
1. 手动输入
2. Jira
3. TM
```

选择 `1` 后，直接在对话框中输入需求描述。

选择 `2` 后，输入 Jira ID，也支持快捷输入：

```text
requirement-case --jiraid bank-123
```

触发后会优先通过已配置的 jira MCP 服务获取 Jira issue 详情，再执行需求分析和自动补全。需要在 Codex 或 Claude 的 MCP 配置中启用 Jira MCP 服务。

选择 `3` 后，输入 TM ID，也支持快捷输入：

```text
requirement-case --tmid tm-123
```

触发后会优先通过已配置的 tm MCP 服务获取 TM 需求详情，再执行需求分析和自动补全。需要在 Codex 或 Claude 的 MCP 配置中启用 TM MCP 服务。

## 阶段缓存与确认

`requirement-case` 端到端流程会在每一步 skill 完成后写入临时 Markdown 缓存：

```text
tmp/.cache/<skill_name>-<需求标识>.md
```

例如：

```text
tmp/.cache/requirement-analysis-bank-123.md
tmp/.cache/functional-testing-bank-123.md
tmp/.cache/testcase-to-excel-bank-123.md
```

每一步完成后会等待用户确认或调整。用户回复 `确认`、`继续` 或 `下一步` 后才执行下一步；回复 `调整：...`、`优化：...` 或 `修改：...` 时，会更新当前阶段缓存后再次等待确认。
