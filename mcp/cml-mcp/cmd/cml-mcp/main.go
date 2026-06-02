package main

import (
	"context"
	"flag"
	"log"
	"log/slog"

	"cml-mcp/internal/config"
	"cml-mcp/internal/logging"
	"cml-mcp/internal/msclient"
	"cml-mcp/internal/server"
)

func main() {
	if err := config.LoadDotEnv(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}
	logger := logging.NewFromEnv()
	slog.SetDefault(logger)
	baseURL := flag.String("base-url", "", "MeterSphere base URL. Defaults to CML_BASE_URL.")
	accessKey := flag.String("access-key", "", "MeterSphere accessKey header. Defaults to CML_ACCESS_KEY.")
	signature := flag.String("signature", "", "MeterSphere signature header. Defaults to CML_SIGNATURE.")
	cookie := flag.String("cookie", "", "Cookie header. Defaults to CML_COOKIE.")
	httpAddr := flag.String("http", "", "HTTP/SSE listen address, such as :8080. Defaults to CML_MCP_HTTP_ADDR or CML_MCP_PORT when set.")
	port := flag.String("port", "", "HTTP/SSE listen port, used when --http is empty.")
	flag.Parse()

	cfg := config.FromEnv()
	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if *accessKey != "" {
		cfg.AccessKey = *accessKey
	}
	if *signature != "" {
		cfg.Signature = *signature
	}
	if *cookie != "" {
		cfg.Cookie = *cookie
	}
	var client *msclient.Client
	if cfg.BaseURL != "" {
		client = msclient.New(cfg, logger)
	}

	addr := *httpAddr
	if addr == "" && *port != "" {
		addr = *port
		if addr[0] != ':' {
			addr = ":" + addr
		}
	}
	if addr == "" {
		addr = config.HTTPAddr()
	}

	var err error
	if addr != "" {
		err = server.RunHTTP(context.Background(), addr, client, logger)
	} else {
		err = server.RunStdio(context.Background(), client, logger)
	}
	if err != nil {
		log.Fatal(err)
	}
}
