package middlewares

import (
	"apiService/internal/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Тесты для RequirePermission

func TestRequirePermission_HasPermission(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("read")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", []string{"read", "write"})
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
}

func TestRequirePermission_MissingPermission(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("admin")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", []string{"read", "write"})
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
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequirePermission_NoPermissionsInContext(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("read")

	router := gin.New()
	// Не устанавливаем permissions в контекст
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequirePermission_InvalidPermissionsFormat(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("read")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", "not-a-slice") // Неправильный тип
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
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequirePermission_EmptyPermissions(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("read")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", []string{})
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
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequirePermission_OptionsRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("read")

	router := gin.New()
	router.Use(middleware)
	router.OPTIONS("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRequirePermission_MultiplePermissions(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("write")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", []string{"read", "write", "delete"})
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
}

func TestRequirePermission_FirstPermissionMatches(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("read")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", []string{"read", "write"})
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
}

func TestRequirePermission_LastPermissionMatches(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	middleware := middlewares.RequirePermission("delete")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("permissions", []string{"read", "write", "delete"})
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
}
