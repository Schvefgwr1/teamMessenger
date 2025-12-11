//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"
)

// logStep - единый helper для читаемых логов в тестах
func logStep(t *testing.T, format string, args ...interface{}) {
	t.Helper()
	t.Logf("==> "+format, args...)
}

// getEnvOrDefault получает значение переменной окружения или возвращает значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
