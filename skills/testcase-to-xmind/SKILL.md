---
name: testcase-to-xmind
description: 将功能测试用例、测试场景、手工 QA 用例、表格用例、Markdown 用例或结构化 JSON 用例转换为 XMind 导图文件。支持用户选择 `1. 通过 XMind 输出`。使用内置 XMind 模板 assets/xmind_case.xmind，按“功能用例 > 模块 > case > 前置条件/备注/标签/步骤描述/文本描述/用例等级”的结构生成导图。适用于用户要求把测试用例转成 xmind、生成 xmind 测试用例导图、导出测试用例脑图或按模板产出 .xmind 文件。
---

# 用例转 XMind

## 概述

使用此 skill 将测试用例整理成 XMind 导图。默认使用 `assets/xmind_case.xmind` 作为模板，保持模板的 XMind 旧版 XML 包格式和节点层级。

当用户提供的是非结构化文本、Markdown 表格、Excel/CSV 表格或普通测试用例描述时，先整理为结构化用例，再生成 XMind。

当用户选择 `1` 或表达“通过 XMind 输出”“通过脑图输出”“通过导图输出”时，使用此 skill。

## 模板结构

模板根节点为：

```text
功能用例
```

推荐生成层级：

```text
功能用例
└── 模块名称
    └── case：用例名称
        ├── 前置条件：...
        ├── 备注：...
        ├── 标签：标签1｜标签2
        ├── 步骤描述
        │   ├── 步骤：...
        │   │   └── 预期结果：...
        │   └── 步骤：...
        │       └── 预期结果：...
        └── 用例等级：P0/P1/P2/P3
```

如果用例没有分步骤，仅有一段操作描述，则使用：

```text
文本描述：...
└── 预期结果：...
```

## 工作流程

1. 解析用户输入，识别模块、用例名称、优先级、前置条件、步骤、预期结果、标签和备注。
2. 如果缺少模块，按功能点、页面、接口或业务流程自动归类；仍无法判断时放入“默认模块”。
3. 如果缺少用例等级，按风险推断：关键主流程用 `P0`，常见流程用 `P1`，边界/异常用 `P2`，低频兼容或展示类检查用 `P3`。
4. 将用例整理为 `references/xmind-case-json-schema.md` 中的 JSON 结构。
5. 使用 `scripts/build_xmind.py` 生成 `.xmind` 文件。
6. 生成后用 `unzip -l` 或脚本输出确认包含 `content.xml`、`styles.xml`、`meta.xml` 和 `META-INF/manifest.xml`。

## 生成命令

把整理后的 JSON 写入临时或工作区文件后运行：

```bash
python3 skills/testcase-to-xmind/scripts/build_xmind.py \
  --input cases.json \
  --output output.xmind
```

默认模板路径为 `skills/testcase-to-xmind/assets/xmind_case.xmind`。如需指定其他模板：

```bash
python3 skills/testcase-to-xmind/scripts/build_xmind.py \
  --input cases.json \
  --template custom_template.xmind \
  --output output.xmind
```

## 输出规则

- 生成的文件扩展名必须是 `.xmind`。
- 用例名称必须有值；缺失时从场景或步骤中提炼。
- 标签多个值使用全角竖线 `｜` 连接，以匹配模板提示。
- 保留中文节点前缀：`case：`、`前置条件：`、`备注：`、`标签：`、`步骤：`、`预期结果：`、`用例等级：`。
- 不确定的字段不要编造；可以省略非必填节点，或写入“待确认”备注。

## 语言

匹配用户提供的用例语言。用户使用中文时，导图节点使用中文。
