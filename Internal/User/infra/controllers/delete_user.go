package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
)

type DeleteUserController struct {
	useCase *app.DeleteUser
}

func NewDeleteUserController(useCase *app.DeleteUser) *DeleteUserController {
	return &DeleteUserController{useCase: useCase}
}

func (ctrl *DeleteUserController) Handle(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de identidad"})
		return
	}

	err := ctrl.useCase.Execute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado permanentemente. Lamentamos que te vayas.",
	})
}