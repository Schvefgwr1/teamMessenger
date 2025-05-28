package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userService/internal/utils"
)

type KeyHandler interface {
	GetPublicKey(c *gin.Context)
}

type keyHandler struct{}

func NewKeyHandler() KeyHandler {
	return &keyHandler{}
}

func (k *keyHandler) GetPublicKey(c *gin.Context) {
	key, err := utils.ExtractPublicKeyFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"key": key})
}
