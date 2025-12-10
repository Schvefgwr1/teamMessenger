package middlewares

import (
	"apiService/internal/middlewares"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis создает тестовый Redis клиент с miniredis
func setupTestRedisForRateLimit(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// Тесты для RateLimitMiddleware

func TestRateLimitMiddleware_WithinLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     10,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
}

func TestRateLimitMiddleware_ExceedsLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     2,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - делаем запросы до превышения лимита
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Третий запрос должен быть заблокирован
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "2", w.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, w.Header().Get("Retry-After"))
}

func TestRateLimitMiddleware_NoUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     10,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	router := gin.New()
	// Не устанавливаем userID в контекст
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - запрос должен пройти, так как нет userID
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitMiddleware_InvalidUserIDType(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     10,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "not-a-uuid") // Неправильный тип
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - запрос должен пройти, так как userID невалидный
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitMiddleware_NilRedis(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	config := middlewares.RateLimitConfig{
		Limit:     10,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(nil, config)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - запрос должен пройти, так как Redis nil (fail open)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitMiddleware_OptionsRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     10,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.Use(middleware)
	router.OPTIONS("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert - OPTIONS запросы должны пропускаться
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRateLimitMiddleware_DifferentUsers(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     2,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	userID1 := uuid.New()
	userID2 := uuid.New()

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - User1 делает 2 запроса (лимит)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.Use(func(c *gin.Context) {
			c.Set("userID", userID1)
			c.Next()
		})
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// User2 должен иметь свой лимит
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID2)
		c.Next()
	})
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitMiddleware_RemainingDecreases(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	config := middlewares.RateLimitConfig{
		Limit:     5,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:test:",
	}

	middleware := middlewares.RateLimitMiddleware(redisClient, config)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - делаем несколько запросов
	var remainingValues []string
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		remainingValues = append(remainingValues, w.Header().Get("X-RateLimit-Remaining"))
	}

	// Assert - remaining должен уменьшаться
	assert.Equal(t, "4", remainingValues[0])
	assert.Equal(t, "3", remainingValues[1])
	assert.Equal(t, "2", remainingValues[2])
}

// Тесты для DefaultAPIRateLimitConfig и StrictRateLimitConfig

func TestDefaultAPIRateLimitConfig(t *testing.T) {
	// Act
	config := middlewares.DefaultAPIRateLimitConfig()

	// Assert
	assert.Equal(t, 200, config.Limit)
	assert.Equal(t, time.Minute, config.Window)
	assert.Equal(t, "ratelimit:user:", config.KeyPrefix)
}

func TestStrictRateLimitConfig(t *testing.T) {
	// Act
	config := middlewares.StrictRateLimitConfig()

	// Assert
	assert.Equal(t, 20, config.Limit)
	assert.Equal(t, time.Minute, config.Window)
	assert.Equal(t, "ratelimit:user:strict:", config.KeyPrefix)
}

// Тесты для RateLimitByEndpoint

func TestRateLimitByEndpoint(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	redisClient := setupTestRedisForRateLimit(t)
	defer redisClient.Close()

	baseConfig := middlewares.RateLimitConfig{
		Limit:     10,
		Window:    time.Minute,
		KeyPrefix: "ratelimit:",
	}

	middleware := middlewares.RateLimitByEndpoint(redisClient, baseConfig, "login")

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.Use(middleware)
	router.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	// Проверяем, что ключ содержит endpoint
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, "ratelimit:login:*").Result()
	require.NoError(t, err)
	assert.Greater(t, len(keys), 0)
}
