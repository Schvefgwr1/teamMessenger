//go:build integration
// +build integration

package integration

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

// createTestFileHeader создает тестовый multipart.FileHeader для использования в тестах
func createTestFileHeader(t *testing.T, filename, contentType string, content []byte) *multipart.FileHeader {
	t.Helper()

	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", filename)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.Write(content)
	if err != nil {
		t.Fatalf("Failed to write content: %v", err)
	}

	// Открываем файл для чтения
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Используем httptest для создания правильного FileHeader
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	file.Close()
	if err != nil {
		t.Fatalf("Failed to copy content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		t.Fatalf("Failed to parse multipart form: %v", err)
	}

	files := req.MultipartForm.File["file"]
	if len(files) == 0 {
		t.Fatalf("No file found in multipart form")
	}

	return files[0]
}
