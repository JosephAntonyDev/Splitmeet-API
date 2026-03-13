package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/app"
	"github.com/gin-gonic/gin"
)

type RegisterDeviceTokenController struct {
	useCase *app.RegisterDeviceToken
}

func NewRegisterDeviceTokenController(useCase *app.RegisterDeviceToken) *RegisterDeviceTokenController {
	return &RegisterDeviceTokenController{useCase: useCase}
}

type registerDeviceTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform"`
}

func (ctrl *RegisterDeviceTokenController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req registerDeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body inválido"})
		return
	}

	input := app.RegisterDeviceTokenInput{
		UserID:   userID.(int64),
		Token:    req.Token,
		Platform: req.Platform,
	}

	if err := ctrl.useCase.Execute(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token de dispositivo registrado"})
}
