---
description: 将结构化测试用例导出为 XMind 或 Excel 文件，支持选项 1/2。
argument-hint: "<1|2|XMind|Excel> <测试用例内容>"
---

请将用户提供的测试用例导出为指定格式。

规则：

- 指定 `1`、XMind、xmind、脑图或导图时，使用 `testcase-to-xmind` skill。
- 指定 `2`、Excel、excel、execl、xlsx 或表格时，使用 `testcase-to-excel` skill。
- 未指定格式时，先向用户展示选项并等待选择：
  1. 通过 XMind 输出
  2. 通过 Excel 输出
- 导出前检查用例字段完整性，包括名称、模块、前置条件、步骤、预期结果、标签和优先级。

完成后输出生成文件路径、用例数量、模板来源和导出过程中的修正说明。
