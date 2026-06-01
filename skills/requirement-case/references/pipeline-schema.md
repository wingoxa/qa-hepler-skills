# 需求到用例数据结构

编排过程中推荐使用以下中间结构：

```json
{
  "requirement_summary": "需求摘要",
  "assumptions": ["明确假设"],
  "cases": [
    {
      "name": "用例名称",
      "module": "/模块/子模块",
      "priority": "P1",
      "precondition": "前置条件",
      "tags": ["需求来源", "功能模块", "场景类型"],
      "remark": "备注",
      "steps": [
        {
          "action": "操作步骤",
          "expected": "预期结果"
        }
      ]
    }
  ],
  "questions": ["待确认问题"]
}
```

## 导入参数

Excel 导入：

```json
{
  "projectId": "项目 ID",
  "file": "/absolute/path/output.xlsx",
  "cover": false,
  "versionId": "可选版本 ID",
  "preCheck": true
}
```

XMind 导入：

```json
{
  "projectId": "项目 ID",
  "file": "/absolute/path/output.xmind",
  "cover": false,
  "versionId": "可选版本 ID",
  "preCheck": true
}
```

## 状态判断

- `preCheck` 成功：继续正式导入。
- `preCheck` 失败：停止导入，输出错误和修复建议。
- 正式导入成功：输出导入数量和平台返回结果。
- 正式导入失败：输出失败原因、文件路径和可重试参数。
