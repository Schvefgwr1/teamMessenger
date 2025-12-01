package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем OPTIONS запросы (preflight) — они обрабатываются CORS middleware
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		perms, exists := c.Get("permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permissions not found in context"})
			return
		}

		permissions, ok := perms.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
			return
		}

		for _, p := range permissions {
			if p == permission {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: missing permission"})
	}
}
