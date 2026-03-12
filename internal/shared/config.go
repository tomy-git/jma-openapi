// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"os"
)

type Config struct {
	Port       string
	Service    string
	Version    string
	JMABaseURL string
	LogFormat  string
}

func LoadConfig() Config {
	return Config{
		Port:       getenv("PORT", "8080"),
		Service:    getenv("SERVICE_NAME", "jma-openapi"),
		Version:    getenv("APP_VERSION", "dev"),
		JMABaseURL: getenv("JMA_BASE_URL", "https://www.jma.go.jp/bosai"),
		LogFormat:  getenv("LOG_FORMAT", "text"),
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}
