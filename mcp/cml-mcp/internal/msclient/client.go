package msclient

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cml-mcp/internal/config"
	"cml-mcp/internal/logging"
	"cml-mcp/internal/mcperr"
)

type Client struct {
	config     config.MeterSphereConfig
	baseURL    *url.URL
	httpClient *http.Client
	logger     *slog.Logger
}

func New(cfg config.MeterSphereConfig, logger *slog.Logger) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if !cfg.VerifySSL {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	if logger != nil {
		logger.Info("metersphere client configured",
			"base_url", cfg.BaseURL,
			"token_configured", cfg.Token != "",
			"cookie_configured", cfg.Cookie != "",
			"verify_ssl", cfg.VerifySSL,
			"timeout_seconds", cfg.Timeout,
		)
	}
	return &Client{
		logger: logger,
		config: cfg,
		httpClient: &http.Client{
			Timeout:   time.Duration(cfg.Timeout * float64(time.Second)),
			Transport: transport,
		},
	}
}

func FromEnv() (*Client, error) {
	cfg := config.FromEnv()
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, mcperr.New("Missing CML_BASE_URL or --base-url for MeterSphere.")
	}
	return New(cfg, logging.NewFromEnv()), nil
}

func (c *Client) endpoint(path string, query map[string]any) (string, error) {
	if c.baseURL == nil {
		parsed, err := url.Parse(strings.TrimRight(c.config.BaseURL, "/") + "/")
		if err != nil {
			return "", err
		}
		c.baseURL = parsed
	}
	relative, err := url.Parse(strings.TrimLeft(path, "/"))
	if err != nil {
		return "", err
	}
	full := c.baseURL.ResolveReference(relative)
	if len(query) > 0 {
		values := full.Query()
		for key, value := range query {
			if value != nil {
				values.Set(key, fmt.Sprint(value))
			}
		}
		full.RawQuery = values.Encode()
	}
	return full.String(), nil
}

func (c *Client) Request(method, path string, payload any, query map[string]any) (any, error) {
	endpoint, err := c.endpoint(path, query)
	if err != nil {
		return nil, err
	}
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(raw)
	}
	request, err := http.NewRequest(strings.ToUpper(method), endpoint, body)
	if err != nil {
		return nil, err
	}
	c.applyHeaders(request)
	if payload != nil {
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	c.logRequest(request, path, payload != nil, 0)
	return c.send(request)
}

func (c *Client) Multipart(path string, requestPayload map[string]any, filePaths []string, fileKey string) (any, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	requestHeader := textproto.MIMEHeader{}
	requestHeader.Set("Content-Disposition", `form-data; name="request"; filename="request.json"`)
	requestHeader.Set("Content-Type", "application/json")
	requestPart, err := writer.CreatePart(requestHeader)
	if err != nil {
		return nil, err
	}
	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}
	if _, err := requestPart.Write(requestJSON); err != nil {
		return nil, err
	}

	for _, filePath := range filePaths {
		if err := addFilePart(writer, fileKey, filePath); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	endpoint, err := c.endpoint(path, nil)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, endpoint, &body)
	if err != nil {
		return nil, err
	}
	c.applyHeaders(request)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	c.logRequest(request, path, true, len(filePaths))
	return c.send(request)
}

func addFilePart(writer *multipart.Writer, fileKey, filePath string) error {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return mcperr.New("File does not exist: " + filePath)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fileKey, filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}

func (c *Client) applyHeaders(request *http.Request) {
	request.Header.Set("Accept", "application/json")
	if c.config.Token != "" {
		token := c.config.Token
		if !strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = "Bearer " + token
		}
		request.Header.Set("Authorization", token)
	}
	if c.config.Cookie != "" {
		request.Header.Set("Cookie", c.config.Cookie)
	}
}

func (c *Client) send(request *http.Request) (any, error) {
	start := time.Now()
	response, err := c.httpClient.Do(request)
	if err != nil {
		c.logResponse(request, 0, start, err)
		return nil, mcperr.New("MeterSphere request failed: " + err.Error())
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		c.logResponse(request, response.StatusCode, start, err)
		return nil, err
	}
	c.logResponse(request, response.StatusCode, start, nil)
	if response.StatusCode >= 400 {
		var parsed any
		if err := json.Unmarshal(raw, &parsed); err != nil {
			parsed = string(raw)
		}
		return nil, mcperr.WithData(fmt.Sprintf("MeterSphere HTTP %d: %s", response.StatusCode, response.Status), parsed)
	}

	contentType := response.Header.Get("Content-Type")
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
		"status":      response.StatusCode,
		"contentType": contentType,
		"bodyBase64":  base64.StdEncoding.EncodeToString(raw),
	}, nil
}

func (c *Client) logRequest(request *http.Request, path string, hasBody bool, fileCount int) {
	if c.logger == nil {
		return
	}
	c.logger.Info("metersphere request", "method", request.Method, "path", path, "body", hasBody, "files", fileCount)
	c.logger.Debug("metersphere request detail", "url", request.URL.Redacted())
}

func (c *Client) logResponse(request *http.Request, status int, start time.Time, err error) {
	if c.logger == nil {
		return
	}
	if err != nil {
		c.logger.Error("metersphere response failed", "method", request.Method, "path", request.URL.Path, "status", status, "duration", logging.Duration(start), "error", err)
		return
	}
	if status >= 400 {
		c.logger.Warn("metersphere response", "method", request.Method, "path", request.URL.Path, "status", status, "duration", logging.Duration(start))
		return
	}
	c.logger.Info("metersphere response", "method", request.Method, "path", request.URL.Path, "status", status, "duration", logging.Duration(start))
}
