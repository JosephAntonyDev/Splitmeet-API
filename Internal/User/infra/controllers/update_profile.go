package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
	"github.com/gin-gonic/gin"
)

type UpdateUserController struct {
	useCase *app.UpdateUser
}

func NewUpdateUserController(useCase *app.UpdateUser) *UpdateUserController {
	return &UpdateUserController{useCase: useCase}
}

type updateUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (ctrl *UpdateUserController) Handle(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	userID := userIDInterface.(int64)

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	updatedUser, err := ctrl.useCase.Execute(app.UpdateUserParams{
		ID:       userID,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
		"user": gin.H{
			"id":       updatedUser.ID,
			"username": updatedUser.Username,
			"name":     updatedUser.Name,
			"email":    updatedUser.Email,
			"phone":    updatedUser.Phone,
		},
	})
}
