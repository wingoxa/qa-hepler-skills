---
description: 分析并自动补全需求，支持来源选项：1 手动输入、2 Jira、3 TM。
argument-hint: "<1|2|3|需求文本|--jiraid JIRA-ID|--tmid TM-ID>"
---

请基于用户输入执行需求分析，优先使用 `requirement-analysis` skill。

如果用户没有直接提供需求描述、Jira ID 或 TM ID，请先展示来源选项并等待选择：

1. 手动输入
2. Jira
3. TM

来源规则：

- 选择 `1`：提示用户直接在对话框中输入需求描述，然后执行分析。
- 选择 `2`：提示用户输入 `jiraid`，调用 jira MCP 服务获取 Jira 需求详情，再执行分析。
- 选择 `3`：提示用户输入 `tmid`，调用 tm MCP 服务获取 TM 需求详情，再执行分析。

如果用户输入包含 `--jiraid <JIRA-ID>`、`jiraid=<JIRA-ID>` 或类似 `requirement-case --jiraid bank-123` 的表达，请先通过已配置的 jira MCP 服务获取 Jira 需求详情，再执行需求分析。

如果用户输入包含 `--tmid <TM-ID>`、`tmid=<TM-ID>` 或类似 `requirement-case --tmid tm-123` 的表达，请先通过已配置的 tm MCP 服务获取 TM 需求详情，再执行需求分析。

Jira 或 TM 获取失败时，不要编造需求内容，说明失败原因并要求用户提供需求文本或修复对应 MCP 配置。

输出内容包括：

1. 需求摘要
2. 来源信息
3. 已知信息
4. 自动补全建议
5. 业务规则
6. 异常与边界
7. 验收标准
8. 问题 / 待确认

请区分明确事实、合理推断和待确认问题，不要把假设写成确定需求。
