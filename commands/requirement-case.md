---
description: 执行需求分析到功能用例设计再到 XMind/Excel 用例文件生成和导入的完整流程。
argument-hint: "<需求文本> [1|2|XMind|Excel] [projectId]"
---

请按 `requirement-case` skill 执行端到端任务流：

1. 需求分析和自动补全
2. 功能测试用例设计
3. XMind 或 Excel 用例文件生成
4. 导入前检查
5. 导入 MeterSphere/CML
6. 结果汇总

生成用例文件前，如果用户没有指定输出格式，必须先展示选项并等待用户选择：

1. 通过 XMind 输出
2. 通过 Excel 输出

用户选择 `1` 时使用 `testcase-to-xmind`；选择 `2`、`excel` 或 `execl` 时使用 `testcase-to-excel`。

每一步 skill 执行完成后，必须将结果写入 `tmp/.cache/<skill_name>-<需求标识>.md`，回复缓存文件路径和阶段摘要，并等待用户确认或调整。用户确认后才能进入下一步；用户要求调整时，更新当前阶段缓存文件并再次等待确认。

如果用户只要求生成文件，不导入平台，则停在文件生成阶段并输出文件路径。

如果缺少项目 ID 且没有 `CML_PROJECT_ID`，不要执行导入，先向用户索要项目 ID。
