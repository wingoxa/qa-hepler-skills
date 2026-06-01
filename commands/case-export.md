---
description: 将结构化测试用例导出为 Excel 或 XMind 文件。
argument-hint: "<Excel|XMind> <测试用例内容>"
---

请将用户提供的测试用例导出为指定格式。

规则：

- 指定 Excel、xlsx 或表格时，使用 `testcase-to-excel` skill。
- 指定 XMind、xmind 或脑图时，使用 `testcase-to-xmind` skill。
- 未指定格式时默认使用 Excel。
- 导出前检查用例字段完整性，包括名称、模块、前置条件、步骤、预期结果、标签和优先级。

完成后输出生成文件路径、用例数量、模板来源和导出过程中的修正说明。
