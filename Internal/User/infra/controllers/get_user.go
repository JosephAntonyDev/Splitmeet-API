package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
)

type GetUserController struct {
	useCase *app.GetUser
}

func NewGetUserController(useCase *app.GetUser) *GetUserController {
	return &GetUserController{useCase: useCase}
}

func (ctrl *GetUserController) Handle(c *gin.Context) {
	idParam := c.Param("id")
	
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	user, err := ctrl.useCase.Execute(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"created_at": user.CreatedAt,
	})
}