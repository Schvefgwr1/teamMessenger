package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userService/internal/services"
	"userService/internal/utils"
)

type KeyHandler interface {
	GetPublicKey(c *gin.Context)
	RegenerateKeys(c *gin.Context)
}

type keyHandler struct {
	keyManagementService *services.KeyManagementService
}

func NewKeyHandler(keyManagementService *services.KeyManagementService) KeyHandler {
	return &keyHandler{
		keyManagementService: keyManagementService,
	}
}

func (k *keyHandler) GetPublicKey(c *gin.Context) {
	key, err := utils.ExtractPublicKeyFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key})
}

// RegenerateKeys ручное обновление ключей
// @Summary Ручное обновление ключей шифрования
// @Description Генерирует новую пару ключей и отправляет публичный ключ в Kafka
// @Tags keys
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Ключи успешно обновлены"
// @Failure 500 {object} map[string]interface{} "Ошибка при обновлении ключей"
// @Router /keys/regenerate [post]
func (k *keyHandler) RegenerateKeys(c *gin.Context) {
	if k.keyManagementService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key management service not available"})
		return
	}

	err := k.keyManagementService.RegenerateKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to regenerate keys: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Keys regenerated successfully",
		"key_version": k.keyManagementService.GetCurrentKeyVersion() - 1, // Текущая версия уже увеличена
	})
}
