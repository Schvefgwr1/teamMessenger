package minio

import (
	"common/config"
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/xerrors"
)

// InitMinIO подключается к MinIO
func InitMinIO(cfg *config.MinIO) (*minio.Client, error) {
	client, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, errors.New("error of creating a minio client")
	}

	if _, err = client.BucketExists(context.Background(), cfg.Bucket); err != nil {
		return nil, xerrors.Errorf("error of connection with MinIO server: %s", err)
	}

	return client, nil
}
