package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type UpdateGroupController struct {
	useCase *app.UpdateGroup
}

func NewUpdateGroupController(useCase *app.UpdateGroup) *UpdateGroupController {
	return &UpdateGroupController{useCase: useCase}
}

type UpdateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ctrl *UpdateGroupController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idParam := c.Param("id")
	groupID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de grupo inválido"})
		return
	}

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	input := app.UpdateGroupInput{
		GroupID:     groupID,
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID.(int64),
	}

	group, err := ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"owner_id":    group.OwnerID,
		"updated_at":  group.UpdatedAt,
	})
}
