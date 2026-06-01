# Excel 用例 JSON 结构

`scripts/build_excel.py` 接收以下 JSON 格式：

```json
{
  "cases": [
    {
      "name": "账号密码登录成功",
      "module": "/登录模块",
      "priority": "P0",
      "precondition": "用户已注册且账号状态正常",
      "remark": "覆盖冒烟主流程",
      "tags": ["登录", "冒烟"],
      "steps": [
        {
          "action": "打开登录页",
          "expected": "展示账号和密码输入框"
        },
        {
          "action": "输入正确账号密码并点击登录",
          "expected": "登录成功并进入首页"
        }
      ]
    }
  ]
}
```

也支持按模块分组：

```json
{
  "modules": [
    {
      "name": "/登录模块",
      "cases": [
        {
          "name": "账号密码登录成功",
          "priority": "P0",
          "steps": [
            {
              "action": "输入正确账号密码并点击登录",
              "expected": "登录成功"
            }
          ]
        }
      ]
    }
  ]
}
```

## 字段说明

| 字段 | 类型 | 必填 | 映射列 | 说明 |
| --- | --- | --- | --- | --- |
| `cases[].name` | string | 是 | 用例名称 | 用例标题 |
| `cases[].module` | string | 否 | 所属模块 | 默认 `/默认模块` |
| `cases[].tags` | array/string | 否 | 标签 | 数组会用 `｜` 拼接 |
| `cases[].precondition` | string | 否 | 前置条件 | 执行前置条件 |
| `cases[].steps` | array | 否 | 步骤描述 / 预期结果 | 分步骤动作和预期 |
| `steps[].action` | string | 否 | 步骤描述 | 单步操作 |
| `steps[].expected` | string | 否 | 预期结果 | 单步预期 |
| `cases[].description` | string | 否 | 步骤描述 | 无分步骤时使用 |
| `cases[].expected` | string | 否 | 预期结果 | 无分步骤时使用 |
| `cases[].edit_mode` | string | 否 | 编辑模式 | 默认 `STEP` |
| `cases[].remark` | string | 否 | 备注 | 补充说明 |
| `cases[].priority` | string | 否 | 用例等级 | `P0`、`P1`、`P2`、`P3` |

## 表格映射建议

从 Markdown、CSV、Excel 或普通表格转换时，优先按以下字段映射：

| 常见表头 | JSON 字段 |
| --- | --- |
| 用例名称、场景、测试场景 | `cases[].name` |
| 所属模块、模块、功能模块 | `cases[].module` |
| 标签、类型、分类 | `cases[].tags` |
| 前置条件、前置 | `cases[].precondition` |
| 步骤、操作步骤、测试步骤、步骤描述 | `cases[].steps[].action` |
| 预期结果、期望结果 | `cases[].steps[].expected` 或 `cases[].expected` |
| 编辑模式、模式 | `cases[].edit_mode` |
| 备注、说明 | `cases[].remark` |
| 优先级、等级、用例等级 | `cases[].priority` |

## 格式规则

- 同一个用例的多步操作拆入 `steps` 数组。
- 如果步骤和预期已经是多行文本，可以直接放入 `description` 和 `expected`。
- 模块缺少 `/` 前缀时，脚本会自动补齐。
- 缺少优先级时默认 `P2`。
