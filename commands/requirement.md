---
description: 分析并自动补全需求，支持通过 --jiraid 从 Jira MCP 获取需求详情。
argument-hint: "<需求文本|--jiraid JIRA-ID>"
---

请基于用户输入执行需求分析，优先使用 `requirement-analysis` skill。

如果用户输入包含 `--jiraid <JIRA-ID>`、`jiraid=<JIRA-ID>` 或类似 `skills --jiraid bank-123` 的表达，请先通过已配置的 jira MCP 服务获取 Jira 需求详情，再执行需求分析。Jira 获取失败时，不要编造需求内容，说明失败原因并要求用户提供需求文本或修复 Jira MCP 配置。

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
