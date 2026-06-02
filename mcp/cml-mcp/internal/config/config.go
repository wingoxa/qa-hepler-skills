package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MeterSphereConfig struct {
	BaseURL   string
	AccessKey string
	Signature string
	Cookie    string
	Timeout   float64
	VerifySSL bool
}

func Env(name, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func LoadDotEnv(paths ...string) error {
	for _, path := range dotEnvPaths(paths...) {
		if err := loadDotEnvFile(path); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
	}
	return nil
}

func dotEnvPaths(paths ...string) []string {
	if len(paths) > 0 {
		return paths
	}
	cwd, err := os.Getwd()
	if err != nil {
		return []string{".env"}
	}
	return []string{
		filepath.Join(cwd, ".env"),
		filepath.Join(cwd, "..", ".env"),
		filepath.Join(cwd, "..", "..", ".env"),
	}
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, ok := parseDotEnvLine(scanner.Text())
		if !ok || os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func parseDotEnvLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
	key, value, ok := strings.Cut(line, "=")
	if !ok {
		return "", "", false
	}
	key = strings.TrimSpace(key)
	if key == "" || strings.ContainsAny(key, " \t") {
		return "", "", false
	}
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return key, value, true
}

func FromEnv() MeterSphereConfig {
	timeout, err := strconv.ParseFloat(Env("CML_TIMEOUT", "60"), 64)
	if err != nil {
		timeout = 60
	}
	verifyValue := strings.ToLower(Env("CML_VERIFY_SSL", "true"))
	return MeterSphereConfig{
		BaseURL:   firstNonEmpty(Env("CML_BASE_URL", ""), Env("METERSPHERE_BASE_URL", "")),
		AccessKey: firstNonEmpty(Env("CML_ACCESS_KEY", ""), Env("METERSPHERE_ACCESS_KEY", "")),
		Signature: firstNonEmpty(Env("CML_SIGNATURE", ""), Env("METERSPHERE_SIGNATURE", "")),
		Cookie:    firstNonEmpty(Env("CML_COOKIE", ""), Env("METERSPHERE_COOKIE", "")),
		Timeout:   timeout,
		VerifySSL: verifyValue != "0" && verifyValue != "false" && verifyValue != "no",
	}
}

func HTTPAddr() string {
	addr := firstNonEmpty(Env("CML_MCP_HTTP_ADDR", ""), Env("CML_MCP_ADDR", ""))
	if addr != "" {
		return addr
	}
	port := firstNonEmpty(Env("CML_MCP_PORT", ""), Env("PORT", ""))
	if port == "" {
		return ""
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func DefaultProjectID() string {
	return firstNonEmpty(Env("CML_PROJECT_ID", ""), Env("METERSPHERE_PROJECT_ID", ""))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
