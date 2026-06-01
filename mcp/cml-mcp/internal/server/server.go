package server

import (
	"context"
	"log/slog"

	"cml-mcp/internal/msclient"
	"cml-mcp/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name    = "cml-mcp"
	Version = "0.1.0"
)

func Run(ctx context.Context, client *msclient.Client, logger *slog.Logger) error {
	if logger != nil {
		logger.Info("starting mcp server", "name", Name, "version", Version, "transport", "stdio")
	}
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: Name, Version: Version}, nil)
	tools.Register(mcpServer, client, logger)
	return mcpServer.Run(ctx, &mcp.StdioTransport{})
}
