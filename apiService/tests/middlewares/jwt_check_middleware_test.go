package middlewares

import (
	"apiService/internal/middlewares"
	"apiService/internal/services"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis создает тестовый Redis клиент с miniredis
func setupTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// generateTestRSAKeyPair генерирует тестовую пару RSA ключей
func generateTestRSAKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

// Claims структура для JWT (должна соответствовать структуре в middleware)
type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

// generateTestJWT создает тестовый JWT токен
func generateTestJWT(t *testing.T, privateKey *rsa.PrivateKey, userID uuid.UUID, permissions []string, expiresAt time.Time) string {
	claims := Claims{
		UserID:      userID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return tokenString
}

// Тесты для JWTMiddlewareWithKeyManager

func TestJWTMiddlewareWithKeyManager_ValidToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	privateKey, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	userID := uuid.New()
	token := generateTestJWT(t, privateKey, userID, []string{"read", "write"}, time.Now().Add(time.Hour))

	// Создаем сессию в Redis
	ctx := context.Background()
	err := sessionService.CreateSession(ctx, userID, token, time.Now().Add(time.Hour))
	require.NoError(t, err)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTMiddlewareWithKeyManager_NoAuthorizationHeader(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddlewareWithKeyManager_InvalidAuthorizationFormat(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddlewareWithKeyManager_NoPublicKey(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	publicKeyManager := services.NewPublicKeyManager() // Ключ не установлен

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJWTMiddlewareWithKeyManager_InvalidToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddlewareWithKeyManager_ExpiredToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	privateKey, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	userID := uuid.New()
	// Токен с истекшим временем
	token := generateTestJWT(t, privateKey, userID, []string{"read"}, time.Now().Add(-time.Hour))

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddlewareWithKeyManager_InvalidSession(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	privateKey, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	userID := uuid.New()
	token := generateTestJWT(t, privateKey, userID, []string{"read"}, time.Now().Add(time.Hour))
	// Не создаем сессию в Redis

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddlewareWithKeyManager_RevokedSession(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	privateKey, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	userID := uuid.New()
	token := generateTestJWT(t, privateKey, userID, []string{"read"}, time.Now().Add(time.Hour))

	// Создаем и затем отзываем сессию
	ctx := context.Background()
	err := sessionService.CreateSession(ctx, userID, token, time.Now().Add(time.Hour))
	require.NoError(t, err)

	err = sessionService.RevokeSession(ctx, userID, token)
	require.NoError(t, err)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddlewareWithKeyManager_OptionsRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.OPTIONS("/protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act - OPTIONS запрос должен пропускаться middleware
	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/protected", nil)
	router.ServeHTTP(w, req)

	// Assert - middleware пропускает OPTIONS, обработчик возвращает 204
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestJWTMiddlewareWithKeyManager_SetsContextValues(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	privateKey, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	userID := uuid.New()
	permissions := []string{"read", "write"}
	token := generateTestJWT(t, privateKey, userID, permissions, time.Now().Add(time.Hour))

	// Создаем сессию в Redis
	ctx := context.Background()
	err := sessionService.CreateSession(ctx, userID, token, time.Now().Add(time.Hour))
	require.NoError(t, err)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	var capturedUserID uuid.UUID
	var capturedPermissions []string
	var capturedToken string
	var capturedKeyVersion int

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		require.True(t, exists)
		capturedUserID = userIDValue.(uuid.UUID)

		permsValue, exists := c.Get("permissions")
		require.True(t, exists)
		capturedPermissions = permsValue.([]string)

		tokenValue, exists := c.Get("token")
		require.True(t, exists)
		capturedToken = tokenValue.(string)

		keyVersionValue, exists := c.Get("keyVersion")
		require.True(t, exists)
		capturedKeyVersion = keyVersionValue.(int)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, permissions, capturedPermissions)
	assert.Equal(t, token, capturedToken)
	assert.Equal(t, 0, capturedKeyVersion) // Начальная версия
}

func TestJWTMiddlewareWithKeyManager_WrongSigningMethod(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyManager := services.NewPublicKeyManager()
	publicKeyManager.SetInitialKey(publicKey)

	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	// Создаем токен с неправильным методом подписи (HS256 вместо RS256)
	claims := jwt.MapClaims{
		"user_id": uuid.New().String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	middleware := middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
