# XMind 用例 JSON 结构

`scripts/build_xmind.py` 接收以下 JSON 格式：

```json
{
  "title": "功能用例",
  "modules": [
    {
      "name": "登录模块",
      "cases": [
        {
          "name": "账号密码登录成功",
          "priority": "P0",
          "precondition": "用户已注册且账号状态正常",
          "remark": "覆盖主流程",
          "tags": ["登录", "冒烟"],
          "steps": [
            {
              "action": "打开登录页，输入正确账号和密码",
              "expected": "登录成功，进入首页"
            }
          ]
        }
      ]
    }
  ]
}
```

## 字段说明

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | string | 否 | 根节点标题，默认 `功能用例` |
| `modules` | array | 是 | 模块列表 |
| `modules[].name` | string | 是 | 模块名称 |
| `modules[].cases` | array | 是 | 模块下的用例列表 |
| `cases[].name` | string | 是 | 用例名称，会生成 `case：用例名称` |
| `cases[].priority` | string | 否 | `P0`、`P1`、`P2`、`P3`，默认 `P2` |
| `cases[].precondition` | string | 否 | 前置条件 |
| `cases[].remark` | string | 否 | 备注 |
| `cases[].tags` | array/string | 否 | 标签，数组会用 `｜` 拼接 |
| `cases[].steps` | array | 否 | 分步骤描述 |
| `steps[].action` | string | 否 | 步骤动作 |
| `steps[].expected` | string | 否 | 预期结果 |
| `cases[].description` | string | 否 | 无分步骤时使用的文本描述 |
| `cases[].expected` | string | 否 | 文本描述下的预期结果 |

## 表格映射建议

从 Markdown、CSV、Excel 或普通表格转换时，优先按以下字段映射：

| 常见表头 | JSON 字段 |
| --- | --- |
| 模块、功能模块、所属模块 | `modules[].name` |
| 用例名称、场景、测试场景 | `cases[].name` |
| 优先级、等级、用例等级 | `cases[].priority` |
| 前置条件、前置 | `cases[].precondition` |
| 步骤、操作步骤、测试步骤 | `cases[].steps[].action` |
| 预期结果、期望结果 | `cases[].steps[].expected` 或 `cases[].expected` |
| 标签、类型、分类 | `cases[].tags` |
| 备注、说明 | `cases[].remark` |

## 归类规则

- 同名模块合并。
- 缺少模块时使用 `默认模块`。
- 缺少优先级时默认 `P2`。
- 多行步骤可以拆成多个 `steps`；如果拆分不可靠，放入 `description`。
