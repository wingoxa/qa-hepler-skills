---
name: requirement-case
description: 将需求分析、功能测试用例设计和测试用例导入串联成完整任务流。适用于用户要求从需求、PRD、用户故事、功能说明或接口说明出发，自动补全需求、设计功能测试用例，并通过选项 1 输出 XMind 或选项 2 输出 Excel 后通过 cml-mcp 导入 MeterSphere/CML。该 agent 会编排 $requirement-analysis、$functional-testing、$testcase-to-excel、$testcase-to-xmind 和 cml-mcp 的 case_import_excel/case_import_xmind 工具。
---

# 需求到用例 Agent

## 目标

使用此 agent 将需求输入转化为可导入测试管理平台的功能测试用例，并执行导入任务。

完整链路：

```text
需求分析 -> 功能用例设计 -> 用例文件生成 -> 导入前检查 -> 用例导入 -> 导入结果汇总
```

每一步 skill 执行完成后，都必须将阶段结果写入临时缓存 Markdown 文件，并等待用户调整或确认后再执行下一步。

## 触发场景

当用户要求以下任务时使用：

- 根据需求生成并导入测试用例
- 从 PRD/用户故事/功能说明创建功能用例并导入 MeterSphere
- 将需求补全后设计测试用例，再导出 Excel 或 XMind
- 执行“需求分析 -> 用例设计 -> 用例导入”的端到端流程

## 工作流程

1. 创建临时缓存目录 `tmp/.cache`。如果目录不存在，先创建。
2. 使用 `$requirement-analysis` 分析需求，补全目标、角色、流程、规则、字段、权限、异常、边界和验收标准。触发需求分析时遵循来源选择：
   ```text
   1. 手动输入
   2. Jira
   3. TM
   ```
   用户选择 `1` 时等待其在对话框输入需求描述；选择 `2` 时要求输入 `jiraid` 并调用 jira-mcp 获取需求详情；选择 `3` 时要求输入 `tmid` 并调用 tm-mcp 获取需求详情。
3. 将需求分析结果写入 `tmp/.cache/requirement-analysis-<需求标识>.md`，回复用户缓存文件路径和阶段摘要，然后等待用户调整或确认。
4. 用户回复确认后，使用 `$functional-testing` 基于确认后的需求分析结果设计功能测试用例。
5. 将功能测试用例结果写入 `tmp/.cache/functional-testing-<需求标识>.md`，回复用户缓存文件路径和阶段摘要，然后等待用户调整或确认。
6. 用户回复确认后，将测试用例规范化为结构化用例数据，至少包含：
   - 用例名称
   - 所属模块
   - 前置条件
   - 步骤描述
   - 预期结果
   - 标签
   - 用例等级
   - 备注
7. 根据用户要求选择导入导出格式：
   - 用户指定 `1`、XMind、xmind、脑图或导图时，使用 `$testcase-to-xmind`。
   - 用户指定 `2`、Excel、excel、execl、xlsx 或表格时，使用 `$testcase-to-excel`。
   - 用户未指定时，必须先展示选项并等待用户选择：
     ```text
     1. 通过 XMind 输出
     2. 通过 Excel 输出
     ```
   - 不要在用户未指定格式时自动默认 Excel。
8. 生成用例文件：
   - Excel：调用 `skills/testcase-to-excel/scripts/build_excel.py`
   - XMind：调用 `skills/testcase-to-xmind/scripts/build_xmind.py`
9. 将用例文件生成结果写入 `tmp/.cache/testcase-to-excel-<需求标识>.md` 或 `tmp/.cache/testcase-to-xmind-<需求标识>.md`，回复用户生成文件路径和缓存文件路径，然后等待用户调整或确认。
10. 用户回复确认导入后，调用 `cml-mcp` 的导入工具：
   - Excel：先调用 `case_import_excel` 且 `preCheck=true`，检查通过后再调用 `case_import_excel` 正式导入。
   - XMind：先调用 `case_import_xmind` 且 `preCheck=true`，检查通过后再调用 `case_import_xmind` 正式导入。
11. 将导入前检查结果写入 `tmp/.cache/case-import-precheck-<需求标识>.md`，回复用户检查结果并等待确认是否正式导入。
12. 用户确认正式导入后执行导入，并将导入结果写入 `tmp/.cache/case-import-<需求标识>.md`。
13. 输出导入结果汇总，包含生成文件路径、缓存文件路径、导入格式、项目 ID、导入数量、失败原因和待确认问题。

## 临时缓存规则

缓存目录固定为：

```text
tmp/.cache
```

缓存文件命名规则：

```text
<skill_name>-<需求标识>.md
```

示例：

```text
tmp/.cache/requirement-analysis-bank-123.md
tmp/.cache/functional-testing-bank-123.md
tmp/.cache/testcase-to-excel-bank-123.md
tmp/.cache/case-import-precheck-bank-123.md
tmp/.cache/case-import-bank-123.md
```

`<需求标识>` 的生成规则：

- 如果用户输入 Jira ID，使用小写 Jira ID，例如 `bank-123`。
- 如果用户指定需求名称，使用需求名称的短横线格式。
- 如果无法提取稳定名称，使用需求摘要前 20 个以内可读字符的短横线格式。
- 文件名仅使用小写字母、数字、短横线和下划线；空格、斜杠、冒号等字符替换为短横线。

每个缓存 Markdown 文件应包含：

1. 阶段名称
2. 需求标识
3. 输入来源
4. 执行时间
5. 本阶段完整结果
6. 待用户确认 / 调整的问题
7. 下一步建议动作

## 用户确认规则

每一步 skill 执行完成后必须暂停，不要自动进入下一步，除非用户已明确要求“全自动执行且无需确认”。

用户可以回复：

- `确认`、`继续`、`下一步`：使用当前缓存结果进入下一步。
- `调整：...`、`优化：...`、`修改：...`：根据用户反馈更新当前阶段结果，并覆盖写回同一个缓存文件，然后再次等待确认。
- `停止`、`暂停`：停止后续动作，并保留已生成的缓存文件。

执行下一步时，必须优先读取上一阶段已确认的缓存文件作为输入，而不是重新从记忆中推断。

## 输入要求

尽量从用户输入中提取：

- 需求文本、PRD、用户故事或接口说明
- 需求描述来源：手动输入、Jira 或 TM
- Jira ID 或 TM ID
- 项目 ID
- 导入格式：Excel 或 XMind
- 是否覆盖已有用例
- 版本 ID
- 模块归属规则

如果缺少项目 ID，应先检查是否已配置 `CML_PROJECT_ID`。如果没有配置，必须向用户索要项目 ID。

## 导入导出选项

用例文件生成、导入前检查和正式导入都使用同一组选项：

1. 通过 XMind 输出
   - 生成 `.xmind` 文件。
   - 使用 `$testcase-to-xmind`。
   - 导入平台时调用 `case_import_xmind`。
2. 通过 Excel 输出
   - 生成 `.xlsx` 文件。
   - 使用 `$testcase-to-excel`。
   - 导入平台时调用 `case_import_excel`。

如果用户输入 `execl`，按 `Excel` 处理。

如果用户未选择，应在功能测试用例确认后、生成文件前暂停并询问用户选择 `1` 或 `2`。

## 用例设计规则

- 用例应覆盖正向流程、反向流程、边界值、权限、状态流转、异常处理和回归风险。
- 用例等级使用 `P0`、`P1`、`P2`、`P3`。
- 主流程和高风险规则优先标为 `P0` 或 `P1`。
- 边界、异常和低频场景通常标为 `P2`。
- 展示类、兼容类或低风险检查可标为 `P3`。
- 所属模块使用平台模板要求的层级格式，例如 `/登录模块`、`/订单/支付`。

## 导入策略

默认执行导入前检查，不直接跳过：

```text
preCheck=true -> 确认无阻塞错误 -> 正式导入
```

如果导入前检查失败：

- 不执行正式导入。
- 汇总失败原因。
- 给出需要修正的用例、字段或模板问题。
- 将失败详情写入 `tmp/.cache/case-import-precheck-<需求标识>.md`。

如果用户明确要求只生成文件，不导入平台：

- 停在文件生成阶段。
- 输出文件路径和导入命令建议。
- 保留生成文件阶段的缓存 Markdown 文件。

## 输出格式

任务完成后输出：

1. 需求分析摘要
2. 用例设计概览
3. 缓存文件
4. 生成文件
5. 导入结果
6. 失败 / 跳过项
7. 待确认问题

## 注意事项

- 不要在缺少项目 ID 且无 `CML_PROJECT_ID` 的情况下执行导入。
- 不要跳过导入前检查，除非用户明确要求。
- 不要把 MCP 认证信息、Token、Cookie 输出到最终回答。
- 如果导入工具不可用，仍应完成需求分析、用例设计和文件生成，并说明导入未执行。
- 缓存文件是临时工作产物，不应写入认证信息。
