package middlewares

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware обрабатывает CORS запросы
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Логируем для отладки (кроме health check)
		if !strings.Contains(c.Request.URL.Path, "health") {
			log.Printf("[CORS] Method=%s Path=%s Origin='%s'", c.Request.Method, c.Request.URL.Path, origin)
		}

		// Проверяем, разрешён ли origin
		if origin != "" && isAllowedOrigin(origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		// Обрабатываем preflight запросы
		if c.Request.Method == "OPTIONS" {
			log.Printf("[CORS] OPTIONS preflight for %s, origin=%s", c.Request.URL.Path, origin)
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin проверяет, разрешён ли origin
func isAllowedOrigin(origin string) bool {
	// Разрешаем localhost на любом порту для разработки
	return strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "http://127.0.0.1:") ||
		origin == "http://localhost" ||
		origin == "http://127.0.0.1"
}
