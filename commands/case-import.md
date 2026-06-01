---
description: 将 Excel 或 XMind 测试用例文件导入 MeterSphere/CML。
argument-hint: "<文件路径> [projectId] [cover=true|false] [versionId]"
---

请使用 `cml-mcp` 将用户提供的 Excel 或 XMind 文件导入 MeterSphere/CML。

规则：

- Excel 文件调用 `case_import_excel`。
- XMind 文件调用 `case_import_xmind`。
- 默认先执行 `preCheck=true`。
- 检查通过后，再执行正式导入。
- 缺少项目 ID 且没有 `CML_PROJECT_ID` 时，先向用户索要项目 ID。
- 不要在输出中暴露 Token、Cookie 或认证头。

完成后输出导入文件、项目 ID、导入数量、失败原因和待确认问题。
