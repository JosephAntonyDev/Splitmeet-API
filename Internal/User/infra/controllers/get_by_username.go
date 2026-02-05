package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
	"github.com/gin-gonic/gin"
)

type GetByUsernameController struct {
	useCase *app.GetByUsername
}

func NewGetByUsernameController(useCase *app.GetByUsername) *GetByUsernameController {
	return &GetByUsernameController{useCase: useCase}
}

func (ctrl *GetByUsernameController) Handle(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username requerido"})
		return
	}

	user, err := ctrl.useCase.Execute(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
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
