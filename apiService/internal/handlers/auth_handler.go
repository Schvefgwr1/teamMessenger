package handlers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	au "common/contracts/api-user"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authController controllers.AuthController
}

func NewAuthHandler(authController controllers.AuthController) *AuthHandler {
	return &AuthHandler{authController: authController}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var registerData dto.RegisterUserRequestGateway

	// Parse JSON part from the multipart form
	jsonData := c.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle optional file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userResponse := h.authController.Register(&registerData, file)
	if userResponse.Error != nil {
		c.JSON(http.StatusInternalServerError, userResponse)
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginData au.Login

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	token, errReq := h.authController.Login(&loginData)
	if errReq != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errReq})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}
