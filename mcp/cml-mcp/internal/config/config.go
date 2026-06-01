package config

import (
	"os"
	"strconv"
	"strings"
)

type MeterSphereConfig struct {
	BaseURL   string
	Token     string
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

func FromEnv() MeterSphereConfig {
	timeout, err := strconv.ParseFloat(Env("CML_TIMEOUT", "60"), 64)
	if err != nil {
		timeout = 60
	}
	verifyValue := strings.ToLower(Env("CML_VERIFY_SSL", "true"))
	return MeterSphereConfig{
		BaseURL:   firstNonEmpty(Env("CML_BASE_URL", ""), Env("METERSPHERE_BASE_URL", "")),
		Token:     firstNonEmpty(Env("CML_TOKEN", ""), Env("METERSPHERE_TOKEN", "")),
		Cookie:    firstNonEmpty(Env("CML_COOKIE", ""), Env("METERSPHERE_COOKIE", "")),
		Timeout:   timeout,
		VerifySSL: verifyValue != "0" && verifyValue != "false" && verifyValue != "no",
	}
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
