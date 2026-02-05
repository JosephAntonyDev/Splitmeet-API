package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type CreateGroupController struct {
	useCase *app.CreateGroup
}

func NewCreateGroupController(useCase *app.CreateGroup) *CreateGroupController {
	return &CreateGroupController{useCase: useCase}
}

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (ctrl *CreateGroupController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: el nombre es requerido"})
		return
	}

	input := app.CreateGroupInput{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID.(int64),
	}

	group, err := ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"owner_id":    group.OwnerID,
		"created_at":  group.CreatedAt,
	})
}
