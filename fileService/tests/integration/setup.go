//go:build integration
// +build integration

package integration

import (
	"common/config"
	"common/db"
	"context"
	"fileService/internal/controllers"
	"fileService/internal/handlers"
	"fileService/internal/repositories"
	"fileService/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testEnv struct {
	Config         *config.Config
	DB             *gorm.DB
	MinioClient    *minio.Client
	FileRepo       repositories.FileRepositoryInterface
	FileTypeRepo   repositories.FileTypeRepositoryInterface
	FileController controllers.FileControllerInterface
}

// setupIntegration настраивает реальные подключения к PostgreSQL и MinIO для интеграционных тестов.
func setupIntegration(t *testing.T) *testEnv {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	cfg := loadTestConfig(t)

	dbClient := connectDatabase(t, cfg)
	require.NoError(t, resetSchema(dbClient))
	require.NoError(t, applyMigrations(dbClient))

	minioClient := initMinio(t, ctx, &cfg.MinIO)

	fileRepo := repositories.NewFileRepository(dbClient)
	fileTypeRepo := repositories.NewFileTypeRepository(dbClient)
	minioAdapter := controllers.NewMinIOAdapter(minioClient)
	fileController := controllers.NewFileController(fileRepo, fileTypeRepo, minioAdapter, &cfg.MinIO)

	env := &testEnv{
		Config:         cfg,
		DB:             dbClient,
		MinioClient:    minioClient,
		FileRepo:       fileRepo,
		FileTypeRepo:   fileTypeRepo,
		FileController: fileController,
	}

	t.Cleanup(func() {
		_ = cleanupBucket(ctx, minioClient, cfg.MinIO.Bucket)
		if sqlDB, err := dbClient.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	return env
}

func loadTestConfig(t *testing.T) *config.Config {
	cfgPath := filepath.Join("..", "..", "config", "config.yaml")
	cfg, err := config.LoadConfig(cfgPath)
	require.NoError(t, err)

	config.ApplyDatabaseEnvOverrides(cfg)
	config.ApplyMinIOEnvOverrides(cfg)
	config.ApplyAppEnvOverrides(cfg)

	// Принудительно устанавливаем тестовые значения для изолированного окружения
	// (перезаписываем значения из конфига, если не заданы через env)
	if os.Getenv("DB_HOST") == "" {
		cfg.Database.Host = "localhost"
	}
	if os.Getenv("DB_USER") == "" {
		cfg.Database.User = "postgres"
	}
	if os.Getenv("DB_PASSWORD") == "" {
		cfg.Database.Password = "postgres"
	}
	if os.Getenv("DB_NAME") == "" {
		cfg.Database.Name = "team_messenger_test"
	}
	if os.Getenv("DB_PORT") == "" {
		cfg.Database.Port = 5433
	}
	if os.Getenv("MINIO_HOST") == "" {
		cfg.MinIO.Host = "localhost:9000"
	}
	if os.Getenv("MINIO_ACCESS_KEY") == "" {
		cfg.MinIO.AccessKey = "minioadmin"
	}
	if os.Getenv("MINIO_SECRET_KEY") == "" {
		cfg.MinIO.SecretKey = "minioadmin"
	}
	if os.Getenv("MINIO_BUCKET") == "" {
		cfg.MinIO.Bucket = "integration-files"
	}

	return cfg
}

func connectDatabase(t *testing.T, cfg *config.Config) *gorm.DB {
	var dbClient *gorm.DB
	var err error

	require.Eventually(t, func() bool {
		dbClient, err = db.InitDB(cfg)
		return err == nil
	}, 10*time.Second, 500*time.Millisecond, "database is not ready: %v", err)

	return dbClient
}

func resetSchema(dbClient *gorm.DB) error {
	resetSQL := "DROP SCHEMA IF EXISTS file_service CASCADE;"
	return dbClient.Exec(resetSQL).Error
}

func applyMigrations(dbClient *gorm.DB) error {
	migrationPath := filepath.Join("..", "..", "migrations", "000001_init_db_v1.up.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := dbClient.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}

func initMinio(t *testing.T, ctx context.Context, cfg *config.MinIO) *minio.Client {
	client, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		_, err = client.ListBuckets(ctx)
		return err == nil
	}, 10*time.Second, 500*time.Millisecond, "minio is not ready: %v", err)

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	require.NoError(t, err)
	if !exists {
		require.NoError(t, client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}))
	}

	return client
}

func cleanupBucket(ctx context.Context, client *minio.Client, bucket string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for object := range client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
			objectsCh <- object
		}
	}()

	for r := range client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}

// newTestRouter собирает gin.Router с реальными зависимостями для интеграционных тестов маршрутов.
func newTestRouter(env *testEnv) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	fileHandler := handlers.NewFileHandler(env.FileController)
	fileTypeController := controllers.NewFileTypeController(env.FileTypeRepo)
	fileTypeHandler := handlers.NewFileTypeHandler(fileTypeController)

	routes.SetupFileRoutes(r, fileHandler)
	routes.SetupFileTypeRoutes(r, fileTypeHandler)
	return r
}
