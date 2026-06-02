package server

import (
	"context"
	"log/slog"
	"net/http"

	"cml-mcp/internal/msclient"
	"cml-mcp/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name    = "cml-mcp"
	Version = "0.1.0"
)

func Run(ctx context.Context, client *msclient.Client, logger *slog.Logger) error {
	return RunStdio(ctx, client, logger)
}

func RunStdio(ctx context.Context, client *msclient.Client, logger *slog.Logger) error {
	if logger != nil {
		logger.Info("starting mcp server", "name", Name, "version", Version, "transport", "stdio")
	}
	mcpServer := New(client, logger)
	return mcpServer.Run(ctx, &mcp.StdioTransport{})
}

func RunHTTP(ctx context.Context, addr string, client *msclient.Client, logger *slog.Logger) error {
	if logger != nil {
		logger.Info("starting mcp server", "name", Name, "version", Version, "transport", "http-sse", "addr", addr)
	}
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: Name, Version: Version}, nil)
	tools.Register(mcpServer, client, logger)
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, nil)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		if err := httpServer.Shutdown(context.Background()); err != nil && logger != nil {
			logger.Warn("http mcp server shutdown failed", "error", err)
		}
	}()
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func New(client *msclient.Client, logger *slog.Logger) *mcp.Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: Name, Version: Version}, nil)
	tools.Register(mcpServer, client, logger)
	return mcpServer
}
