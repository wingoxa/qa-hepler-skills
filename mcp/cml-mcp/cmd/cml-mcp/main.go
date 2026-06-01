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
	logger := logging.NewFromEnv()
	slog.SetDefault(logger)
	baseURL := flag.String("base-url", "", "MeterSphere base URL. Defaults to CML_BASE_URL.")
	token := flag.String("token", "", "Bearer token. Defaults to CML_TOKEN.")
	cookie := flag.String("cookie", "", "Cookie header. Defaults to CML_COOKIE.")
	flag.Parse()

	var client *msclient.Client
	if *baseURL != "" {
		cfg := config.FromEnv()
		cfg.BaseURL = *baseURL
		if *token != "" {
			cfg.Token = *token
		}
		if *cookie != "" {
			cfg.Cookie = *cookie
		}
		client = msclient.New(cfg, logger)
	}

	if err := server.Run(context.Background(), client, logger); err != nil {
		log.Fatal(err)
	}
}
