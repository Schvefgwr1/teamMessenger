//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"testing"
	"time"
)

func TestFileController_UploadAndGet_Integration(t *testing.T) {
	env := setupIntegration(t)
	ctx := context.Background()

	fileHeader, cleanup, err := createMultipartFile("sample.png", "image/png", []byte("hello from integration"))
	require.NoError(t, err)
	defer cleanup()

	uploaded, err := env.FileController.UploadFile(fileHeader)
	require.NoError(t, err)
	require.NotNil(t, uploaded)

	objectInfo, err := env.MinioClient.StatObject(ctx, env.Config.MinIO.Bucket, uploaded.Name, minio.StatObjectOptions{})
	require.NoError(t, err)
	assert.Equal(t, fileHeader.Size, objectInfo.Size)

	stored, err := env.FileController.GetFile(uploaded.ID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, uploaded.ID, stored.ID)
	assert.Equal(t, uploaded.Name, stored.Name)
	assert.Contains(t, stored.URL, env.Config.MinIO.Bucket)
}

func TestFileController_RenameFile_Integration(t *testing.T) {
	env := setupIntegration(t)
	ctx := context.Background()

	fileHeader, cleanup, err := createMultipartFile("to-rename.png", "image/png", []byte("rename me"))
	require.NoError(t, err)
	defer cleanup()

	uploaded, err := env.FileController.UploadFile(fileHeader)
	require.NoError(t, err)

	oldName := uploaded.Name
	newName := fmt.Sprintf("renamed-%d", time.Now().UnixNano())

	updated, err := env.FileController.RenameFile(uploaded.ID, newName)
	require.NoError(t, err)
	require.NotNil(t, updated)

	_, err = env.MinioClient.StatObject(ctx, env.Config.MinIO.Bucket, updated.Name, minio.StatObjectOptions{})
	require.NoError(t, err)

	_, err = env.MinioClient.StatObject(ctx, env.Config.MinIO.Bucket, oldName, minio.StatObjectOptions{})
	require.Error(t, err)

	assert.Contains(t, updated.URL, updated.Name)
}

func createMultipartFile(filename, contentType string, content []byte) (*multipart.FileHeader, func(), error) {
	if len(content) == 0 {
		content = []byte("integration content")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, func() {}, err
	}

	if _, err = part.Write(content); err != nil {
		return nil, func() {}, err
	}

	if err = writer.Close(); err != nil {
		return nil, func() {}, err
	}

	reader := multipart.NewReader(&body, writer.Boundary())
	form, err := reader.ReadForm(32 << 20)
	if err != nil {
		return nil, func() {}, err
	}

	fileHeader := form.File["file"][0]
	fileHeader.Header.Set("Content-Type", contentType)

	cleanup := func() {
		form.RemoveAll()
	}

	return fileHeader, cleanup, nil
}
