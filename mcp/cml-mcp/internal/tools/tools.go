package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"cml-mcp/internal/config"
	"cml-mcp/internal/logging"
	"cml-mcp/internal/msclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Definition map[string]any

func Register(server *mcp.Server, client *msclient.Client, logger *slog.Logger) {
	for _, definition := range Schema() {
		name := fmt.Sprint(definition["name"])
		description := fmt.Sprint(definition["description"])
		inputSchema := definition["inputSchema"]
		server.AddTool(&mcp.Tool{
			Name:        name,
			Description: description,
			InputSchema: inputSchema,
		}, handler(name, client, logger))
		if logger != nil {
			logger.Debug("registered tool", "name", name)
		}
	}
}

func handler(name string, configuredClient *msclient.Client, logger *slog.Logger) mcp.ToolHandler {
	return func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		args := map[string]any{}
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				logToolError(logger, name, start, err)
				return toolError("invalid arguments: " + err.Error()), nil
			}
		}
		if logger != nil {
			logger.Info("tool call start", "name", name, "arg_keys", logging.Keys(args))
		}

		client := configuredClient
		if client == nil {
			var err error
			client, err = msclient.FromEnv()
			if err != nil {
				logToolError(logger, name, start, err)
				return toolError(err.Error()), nil
			}
		}

		data, err := Call(name, args, client)
		if err != nil {
			logToolError(logger, name, start, err)
			return toolError(err.Error()), nil
		}
		pretty, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			logToolError(logger, name, start, err)
			return toolError(err.Error()), nil
		}
		if logger != nil {
			logger.Info("tool call done", "name", name, "duration", logging.Duration(start))
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(pretty)}},
		}, nil
	}
}

func logToolError(logger *slog.Logger, name string, start time.Time, err error) {
	if logger != nil {
		logger.Error("tool call failed", "name", name, "duration", logging.Duration(start), "error", err)
	}
}

func toolError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
		IsError: true,
	}
}

func Schema() []Definition {
	authNote := "需要通过环境变量 CML_BASE_URL + CML_ACCESS_KEY + CML_SIGNATURE 或 CML_COOKIE 配置 MeterSphere 访问。"
	return []Definition{
		{
			"name":        "case_query",
			"description": "查询 MeterSphere 功能测试用例列表。对应 POST /functional/case/page。" + authNote,
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"projectId": map[string]any{"type": "string", "description": "项目 ID。未传时使用 CML_PROJECT_ID。"},
					"current":   map[string]any{"type": "integer", "default": 1},
					"pageSize":  map[string]any{"type": "integer", "default": 20},
					"moduleIds": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"keyword":   map[string]any{"type": "string", "description": "透传给 MeterSphere 的搜索关键字字段。"},
					"request":   map[string]any{"type": "object", "description": "额外透传或覆盖的 FunctionalCasePageRequest 字段。"},
				},
			},
		},
		{
			"name":        "case_add",
			"description": "新增 MeterSphere 功能测试用例。对应 multipart POST /functional/case/add。",
			"inputSchema": objectSchema(map[string]any{
				"request": map[string]any{"type": "object", "description": "FunctionalCaseAddRequest 原始请求体。"},
				"files":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "可选附件路径列表。"},
			}, []string{"request"}),
		},
		{
			"name":        "case_detail",
			"description": "查看 MeterSphere 功能测试用例详情。对应 GET /functional/case/detail/{id}。",
			"inputSchema": objectSchema(map[string]any{
				"id": map[string]any{"type": "string", "description": "功能用例 ID。"},
			}, []string{"id"}),
		},
		{
			"name":        "case_update",
			"description": "更新 MeterSphere 功能测试用例详情。对应 multipart POST /functional/case/update。",
			"inputSchema": objectSchema(map[string]any{
				"request": map[string]any{"type": "object", "description": "FunctionalCaseEditRequest 原始请求体，必须包含 id。"},
				"files":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "可选附件路径列表。"},
			}, []string{"request"}),
		},
		{
			"name":        "case_import_xmind",
			"description": "导入 XMind 测试用例文件。可先 preCheck，再正式导入。对应 /functional/case/pre-check/xmind 和 /functional/case/import/xmind。",
			"inputSchema": importSchema(".xmind"),
		},
		{
			"name":        "case_import_excel",
			"description": "导入 Excel 测试用例文件。可先 preCheck，再正式导入。对应 /functional/case/pre-check/excel 和 /functional/case/import/excel。",
			"inputSchema": importSchema(".xlsx"),
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	return map[string]any{"type": "object", "properties": properties, "required": required}
}

func importSchema(ext string) map[string]any {
	return objectSchema(map[string]any{
		"projectId": map[string]any{"type": "string", "description": "项目 ID。未传时使用 CML_PROJECT_ID。"},
		"file":      map[string]any{"type": "string", "description": "本地 " + ext + " 文件路径。"},
		"cover":     map[string]any{"type": "boolean", "default": false},
		"versionId": map[string]any{"type": "string"},
		"count":     map[string]any{"type": "string"},
		"preCheck":  map[string]any{"type": "boolean", "default": false},
		"request":   map[string]any{"type": "object", "description": "额外透传或覆盖的 FunctionalCaseImportRequest 字段。"},
	}, []string{"file"})
}

func Call(name string, args map[string]any, client *msclient.Client) (any, error) {
	switch name {
	case "case_query":
		request := mapValue(args["request"])
		if err := ensureProjectID(request, args); err != nil {
			return nil, err
		}
		setDefault(request, "current", intValue(args["current"], 1))
		setDefault(request, "pageSize", intValue(args["pageSize"], 20))
		if value, ok := args["moduleIds"]; ok {
			request["moduleIds"] = value
		}
		if value, ok := args["keyword"]; ok && fmt.Sprint(value) != "" {
			request["keyword"] = value
		}
		return client.Request("POST", "/functional/case/page", request, nil)
	case "case_add":
		request, err := requiredMap(args, "request")
		if err != nil {
			return nil, err
		}
		return client.Multipart("/functional/case/add", request, stringSlice(args["files"]), "files")
	case "case_detail":
		id := fmt.Sprint(args["id"])
		if id == "" {
			return nil, fmt.Errorf("missing required argument: id")
		}
		return client.Request("GET", "/functional/case/detail/"+id, nil, nil)
	case "case_update":
		request, err := requiredMap(args, "request")
		if err != nil {
			return nil, err
		}
		return client.Multipart("/functional/case/update", request, stringSlice(args["files"]), "files")
	case "case_import_xmind":
		request, err := importRequest(args)
		if err != nil {
			return nil, err
		}
		path := "/functional/case/import/xmind"
		if boolValue(args["preCheck"]) {
			path = "/functional/case/pre-check/xmind"
		}
		return client.Multipart(path, request, []string{fmt.Sprint(args["file"])}, "file")
	case "case_import_excel":
		request, err := importRequest(args)
		if err != nil {
			return nil, err
		}
		path := "/functional/case/import/excel"
		if boolValue(args["preCheck"]) {
			path = "/functional/case/pre-check/excel"
		}
		return client.Multipart(path, request, []string{fmt.Sprint(args["file"])}, "file")
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func projectID(args map[string]any) (string, error) {
	if value, ok := args["projectId"]; ok && fmt.Sprint(value) != "" {
		return fmt.Sprint(value), nil
	}
	value := config.DefaultProjectID()
	if value == "" {
		return "", fmt.Errorf("missing projectId. Pass projectId or set CML_PROJECT_ID")
	}
	return value, nil
}

func ensureProjectID(payload map[string]any, args map[string]any) error {
	if value, ok := payload["projectId"]; ok && fmt.Sprint(value) != "" {
		return nil
	}
	projectID, err := projectID(args)
	if err != nil {
		return err
	}
	payload["projectId"] = projectID
	return nil
}

func importRequest(args map[string]any) (map[string]any, error) {
	file := fmt.Sprint(args["file"])
	if file == "" {
		return nil, fmt.Errorf("missing required argument: file")
	}
	payload := mapValue(args["request"])
	if err := ensureProjectID(payload, args); err != nil {
		return nil, err
	}
	setDefault(payload, "cover", boolValue(args["cover"]))
	if value, ok := args["versionId"]; ok {
		payload["versionId"] = value
	}
	if value, ok := args["count"]; ok {
		payload["count"] = fmt.Sprint(value)
	}
	return payload, nil
}

func requiredMap(args map[string]any, key string) (map[string]any, error) {
	value, ok := args[key]
	if !ok {
		return nil, fmt.Errorf("missing required argument: %s", key)
	}
	return mapValue(value), nil
}

func mapValue(value any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	if mapped, ok := value.(map[string]any); ok {
		return mapped
	}
	return map[string]any{}
}

func setDefault(target map[string]any, key string, value any) {
	if _, ok := target[key]; !ok {
		target[key] = value
	}
}

func stringSlice(value any) []string {
	if value == nil {
		return nil
	}
	if items, ok := value.([]any); ok {
		result := make([]string, 0, len(items))
		for _, item := range items {
			result = append(result, fmt.Sprint(item))
		}
		return result
	}
	if items, ok := value.([]string); ok {
		return items
	}
	return nil
}

func boolValue(value any) bool {
	if value == nil {
		return false
	}
	if boolean, ok := value.(bool); ok {
		return boolean
	}
	parsed, _ := strconv.ParseBool(fmt.Sprint(value))
	return parsed
}

func intValue(value any, fallback int) int {
	if value == nil {
		return fallback
	}
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	default:
		parsed, err := strconv.Atoi(fmt.Sprint(value))
		if err != nil {
			return fallback
		}
		return parsed
	}
}
