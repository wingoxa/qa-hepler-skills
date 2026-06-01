# cml-mcp

`cml-mcp` 是一个 stdio MCP 服务，用于调用 MeterSphere 功能用例接口。

## 功能

- `case_query`：测试用例查询，调用 `POST /functional/case/page`
- `case_add`：测试用例添加，调用 `POST /functional/case/add`
- `case_detail`：查看用例详情，调用 `GET /functional/case/detail/{id}`
- `case_update`：更新用例详情，调用 `POST /functional/case/update`
- `case_import_xmind`：导入 XMind 用例，调用 `POST /functional/case/import/xmind`
- `case_import_excel`：导入 Excel 用例，调用 `POST /functional/case/import/excel`

## 配置

通过环境变量配置 MeterSphere：

```bash
export CML_BASE_URL="http://localhost:8081"
export CML_TOKEN="your-token"
export CML_PROJECT_ID="project-id"
```

如果实例使用 Cookie 认证：

```bash
export CML_COOKIE="SESSION=..."
```

日志默认输出到 `stderr`，不会污染 MCP 的 `stdout` 协议流。可通过 `CML_LOG_LEVEL` 控制级别：

```bash
export CML_LOG_LEVEL="debug" # debug, info, warn, error, off
```

## 启动

```bash
cd mcp/cml-mcp
go run ./cmd/cml-mcp
```

也可以先构建二进制：

```bash
cd mcp/cml-mcp
go build -o cml-mcp ./cmd/cml-mcp
./cml-mcp
```

MCP 客户端配置示例：

```json
{
  "mcpServers": {
    "cml-mcp": {
      "command": "go",
      "args": ["run", "./cmd/cml-mcp"],
      "cwd": "/Users/black/Desktop/codings/qa-hepler-skills/mcp/cml-mcp",
      "env": {
        "CML_BASE_URL": "http://localhost:8081",
        "CML_TOKEN": "your-token",
        "CML_PROJECT_ID": "project-id",
        "CML_LOG_LEVEL": "info"
      }
    }
  }
}
```

## 代码结构

```text
cmd/cml-mcp/
└── main.go              # CLI 参数解析和服务启动
internal/
├── config/config.go     # 环境变量和配置对象
├── logging/logging.go   # stderr 日志和日志级别
├── mcperr/error.go      # MCP 错误类型
├── msclient/client.go   # MeterSphere HTTP / multipart 客户端
├── server/server.go     # go-sdk MCP Server 和 stdio transport
└── tools/tools.go       # MCP 工具 schema 和工具调用分发
```

## 请求体说明

新增和更新用例会透传 MeterSphere 的 `FunctionalCaseAddRequest` / `FunctionalCaseEditRequest`，并用 multipart 的 `request` 字段提交 JSON。

导入用例会提交 `FunctionalCaseImportRequest`：

```json
{
  "projectId": "project-id",
  "cover": false,
  "versionId": "optional",
  "count": "optional"
}
```
