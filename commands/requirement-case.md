---
description: 执行需求分析到功能用例设计再到用例文件生成和导入的完整流程。
argument-hint: "<需求文本> [Excel|XMind] [projectId]"
---

请按 `requirement-case` skill 执行端到端任务流：

1. 需求分析和自动补全
2. 功能测试用例设计
3. Excel 或 XMind 用例文件生成
4. 导入前检查
5. 导入 MeterSphere/CML
6. 结果汇总

每一步 skill 执行完成后，必须将结果写入 `tmp/.cache/<skill_name>-<需求标识>.md`，回复缓存文件路径和阶段摘要，并等待用户确认或调整。用户确认后才能进入下一步；用户要求调整时，更新当前阶段缓存文件并再次等待确认。

如果用户只要求生成文件，不导入平台，则停在文件生成阶段并输出文件路径。

如果缺少项目 ID 且没有 `CML_PROJECT_ID`，不要执行导入，先向用户索要项目 ID。
