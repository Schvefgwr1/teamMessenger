package middlewares

import (
	"apiService/internal/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

// JWTMiddlewareWithKeyManager создает middleware с динамическим обновлением ключей
func JWTMiddlewareWithKeyManager(publicKeyManager *services.PublicKeyManager, sessionService *services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем OPTIONS запросы (preflight) — они обрабатываются CORS middleware
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		type Claims struct {
			UserID      uuid.UUID `json:"user_id"`
			Permissions []string  `json:"permissions"`
			jwt.RegisteredClaims
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Получаем текущий публичный ключ (thread-safe)
		publicKey := publicKeyManager.GetCurrentKey()
		if publicKey == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Public key not available"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Проверяем сессию в Redis
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		isValid, err := sessionService.IsSessionValid(ctx, claims.UserID, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session validation failed"})
			return
		}

		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is invalid or revoked"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("permissions", claims.Permissions)
		c.Set("token", tokenStr)
		c.Set("keyVersion", publicKeyManager.GetKeyVersion()) // Добавляем версию ключа для отладки
		c.Next()
	}
}
