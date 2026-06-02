# cml-mcp

`cml-mcp` 是一个 stdio / HTTP SSE MCP 服务，用于调用 MeterSphere 功能用例接口。

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
export CML_ACCESS_KEY="your-access-key"
export CML_SIGNATURE="your-signature"
export CML_PROJECT_ID="project-id"
```

也可以在当前目录、上级目录或项目根目录放置 `.env` 文件。程序启动时会自动读取 `.env`，但不会覆盖已经存在的系统环境变量或命令行参数。可参考：

```text
mcp/cml-mcp/.env.example
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

HTTP SSE / streamable HTTP 模式：

```bash
cd mcp/cml-mcp
go run ./cmd/cml-mcp --http :8080
```

也可以通过端口或环境变量启动：

```bash
go run ./cmd/cml-mcp --port 8080
export CML_MCP_HTTP_ADDR=":8080"
export CML_MCP_PORT="8080"
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
        "CML_ACCESS_KEY": "your-access-key",
        "CML_SIGNATURE": "your-signature",
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
├── msclient/client.go   # MeterSphere resty HTTP / multipart 客户端
├── server/server.go     # go-sdk MCP Server、stdio 和 HTTP SSE transport
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
