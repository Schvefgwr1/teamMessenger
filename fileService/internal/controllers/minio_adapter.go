package controllers

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

// MinIOAdapter адаптирует *minio.Client к интерфейсу MinIOClientInterface
type MinIOAdapter struct {
	client *minio.Client
}

// NewMinIOAdapter создает новый адаптер для MinIO клиента
func NewMinIOAdapter(client *minio.Client) *MinIOAdapter {
	return &MinIOAdapter{client: client}
}

func (a *MinIOAdapter) StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	return a.client.StatObject(ctx, bucketName, objectName, opts)
}

func (a *MinIOAdapter) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return a.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (a *MinIOAdapter) CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error) {
	return a.client.CopyObject(ctx, dst, src)
}

func (a *MinIOAdapter) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	return a.client.RemoveObject(ctx, bucketName, objectName, opts)
}
