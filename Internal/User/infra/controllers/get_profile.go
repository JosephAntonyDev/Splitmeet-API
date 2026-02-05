package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
	"github.com/gin-gonic/gin"
)

type GetProfileController struct {
	useCase *app.GetProfile
}

func NewGetProfileController(useCase *app.GetProfile) *GetProfileController {
	return &GetProfileController{useCase: useCase}
}

func (ctrl *GetProfileController) Handle(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token no procesado"})
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno procesando identidad del usuario"})
		return
	}

	user, err := ctrl.useCase.Execute(userID)

	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"created_at": user.CreatedAt,
	})
}
