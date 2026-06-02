package msclient

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"cml-mcp/internal/config"
	"cml-mcp/internal/logging"
	"cml-mcp/internal/mcperr"

	"resty.dev/v3"
)

type Client struct {
	config     config.MeterSphereConfig
	httpClient *resty.Client
	logger     *slog.Logger
}

func New(cfg config.MeterSphereConfig, logger *slog.Logger) *Client {
	httpClient := resty.New().
		SetBaseURL(strings.TrimRight(cfg.BaseURL, "/")).
		SetTimeout(time.Duration(cfg.Timeout*float64(time.Second))).
		SetHeader("Accept", "application/json")
	if !cfg.VerifySSL {
		httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}
	if logger != nil {
		logger.Info("metersphere client configured",
			"base_url", cfg.BaseURL,
			"access_key_configured", cfg.AccessKey != "",
			"signature_configured", cfg.Signature != "",
			"cookie_configured", cfg.Cookie != "",
			"verify_ssl", cfg.VerifySSL,
			"timeout_seconds", cfg.Timeout,
		)
	}
	return &Client{
		logger:     logger,
		config:     cfg,
		httpClient: httpClient,
	}
}

func FromEnv() (*Client, error) {
	cfg := config.FromEnv()
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, mcperr.New("Missing CML_BASE_URL or --base-url for MeterSphere.")
	}
	return New(cfg, logging.NewFromEnv()), nil
}

func (c *Client) Request(method, path string, payload any, query map[string]any) (any, error) {
	request := c.request()
	c.applyQuery(request, query)
	if payload != nil {
		request.SetHeader("Content-Type", "application/json; charset=utf-8")
		request.SetBody(payload)
	}
	c.logRequest(strings.ToUpper(method), path, payload != nil, 0)
	return c.send(strings.ToUpper(method), path, request)
}

func (c *Client) Multipart(path string, requestPayload map[string]any, filePaths []string, fileKey string) (any, error) {
	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	request := c.request().
		SetMultipartField("request", "request.json", "application/json", bytes.NewReader(requestJSON))
	for _, filePath := range filePaths {
		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}
		request.SetFile(fileKey, filePath)
	}
	c.logRequest("POST", path, true, len(filePaths))
	return c.send("POST", path, request)
}

func (c *Client) request() *resty.Request {
	request := c.httpClient.R().
		SetHeader("Accept", "application/json")
	if c.config.AccessKey != "" {
		request.SetHeader("accessKey", c.config.AccessKey)
	}
	if c.config.Signature != "" {
		request.SetHeader("signature", c.config.Signature)
	}
	if c.config.Cookie != "" {
		request.SetHeader("Cookie", c.config.Cookie)
	}
	return request
}

func (c *Client) applyQuery(request *resty.Request, query map[string]any) {
	for key, value := range query {
		if value != nil {
			request.SetQueryParam(key, fmt.Sprint(value))
		}
	}
}

func (c *Client) send(method, path string, request *resty.Request) (any, error) {
	start := time.Now()
	response, err := request.Execute(method, path)
	if err != nil {
		c.logResponse(method, path, 0, start, err)
		return nil, mcperr.New("MeterSphere request failed: " + err.Error())
	}

	raw := response.Bytes()
	c.logResponse(method, path, response.StatusCode(), start, nil)
	if response.StatusCode() >= 400 {
		var parsed any
		if err := json.Unmarshal(raw, &parsed); err != nil {
			parsed = string(raw)
		}
		return nil, mcperr.WithData(fmt.Sprintf("MeterSphere HTTP %d: %s", response.StatusCode(), response.Status()), parsed)
	}

	contentType := response.Header().Get("Content-Type")
	if strings.Contains(contentType, "application/json") || bytes.HasPrefix(bytes.TrimSpace(raw), []byte("{")) || bytes.HasPrefix(bytes.TrimSpace(raw), []byte("[")) {
		var parsed any
		if len(raw) == 0 {
			return nil, nil
		}
		if err := json.Unmarshal(raw, &parsed); err != nil {
			return nil, err
		}
		return parsed, nil
	}
	return map[string]any{
		"status":      response.StatusCode(),
		"contentType": contentType,
		"bodyBase64":  base64.StdEncoding.EncodeToString(raw),
	}, nil
}

func (c *Client) logRequest(method, path string, hasBody bool, fileCount int) {
	if c.logger == nil {
		return
	}
	c.logger.Info("metersphere request", "method", method, "path", path, "body", hasBody, "files", fileCount)
}

func (c *Client) logResponse(method, path string, status int, start time.Time, err error) {
	if c.logger == nil {
		return
	}
	if err != nil {
		c.logger.Error("metersphere response failed", "method", method, "path", path, "status", status, "duration", logging.Duration(start), "error", err)
		return
	}
	if status >= 400 {
		c.logger.Warn("metersphere response", "method", method, "path", path, "status", status, "duration", logging.Duration(start))
		return
	}
	c.logger.Info("metersphere response", "method", method, "path", path, "status", status, "duration", logging.Duration(start))
}
