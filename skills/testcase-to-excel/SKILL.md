---
name: testcase-to-excel
description: 将功能测试用例、测试场景、手工 QA 用例、Markdown 表格、CSV/JSON 用例或普通文本用例转换为 Excel 测试用例文件。支持用户选择 `2. 通过 Excel 输出`，并兼容 `execl` 拼写。使用内置 Excel 模板 assets/excel_case.xlsx，按“用例名称、所属模块、标签、前置条件、步骤描述、预期结果、编辑模式、备注、用例等级”的列结构生成 .xlsx。适用于用户要求把测试用例转成 excel、生成测试用例表格、导出用例 xlsx、按模板产出 Excel 测试用例文件。
---

# 用例转 Excel

## 概述

使用此 skill 将测试用例整理成 Excel 文件。默认使用 `assets/excel_case.xlsx` 作为模板，保留模板 sheet、表头和基础样式。

当用户提供的是非结构化文本、Markdown 表格、XMind 结构、普通测试场景或功能说明时，先整理为结构化用例，再生成 Excel。

当用户选择 `2` 或表达“通过 Excel 输出”“通过 execl 输出”“通过表格输出”时，使用此 skill。

## 模板结构

模板文件：`assets/excel_case.xlsx`

默认工作表：`模版`

表头顺序：

| 列 | 字段 |
| --- | --- |
| A | 用例名称 |
| B | 所属模块 |
| C | 标签 |
| D | 前置条件 |
| E | 步骤描述 |
| F | 预期结果 |
| G | 编辑模式 |
| H | 备注 |
| I | 用例等级 |

数据从第 2 行开始。模板示例中 `编辑模式` 使用 `STEP`，步骤和预期结果使用编号格式：

```text
[1]打开登录页
[2]输入账号密码
```

## 工作流程

1. 解析用户输入，识别用例名称、所属模块、标签、前置条件、步骤、预期结果、备注和用例等级。
2. 如果缺少所属模块，按页面、功能、接口或业务流程自动归类；仍无法判断时使用 `/默认模块`。
3. 如果模块有层级，使用 `/一级模块/二级模块` 格式。
4. 如果缺少用例等级，按风险推断：关键主流程用 `P0`，常见流程用 `P1`，边界/异常用 `P2`，低频兼容或展示类检查用 `P3`。
5. 将用例整理为 `references/excel-case-json-schema.md` 中的 JSON 结构。
6. 使用 `scripts/build_excel.py` 生成 `.xlsx` 文件。
7. 生成后读取工作簿，确认 sheet、表头和数据行数正确。

## 生成命令

把整理后的 JSON 写入临时或工作区文件后运行：

```bash
python3 skills/testcase-to-excel/scripts/build_excel.py \
  --input cases.json \
  --output output.xlsx
```

默认模板路径为 `skills/testcase-to-excel/assets/excel_case.xlsx`。如需指定其他模板：

```bash
python3 skills/testcase-to-excel/scripts/build_excel.py \
  --input cases.json \
  --template custom_template.xlsx \
  --output output.xlsx
```

## 输出规则

- 生成的文件扩展名必须是 `.xlsx`。
- 用例名称必须有值；缺失时从场景或步骤中提炼。
- 所属模块必须以 `/` 开头，例如 `/登录模块` 或 `/订单/支付`。
- 标签多个值使用英文逗号或全角竖线均可；脚本默认用 `｜` 连接数组标签。
- 步骤和预期结果逐行编号，格式为 `[1]...`、`[2]...`。
- `编辑模式` 默认填 `STEP`。
- `用例等级` 只使用 `P0`、`P1`、`P2`、`P3`。
- 不确定的非必填字段可以留空；不要编造业务规则。

## 语言

匹配用户提供的用例语言。用户使用中文时，Excel 内容使用中文。
