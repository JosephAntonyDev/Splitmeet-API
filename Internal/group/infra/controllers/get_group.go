package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type GetGroupController struct {
	useCase *app.GetGroup
}

func NewGetGroupController(useCase *app.GetGroup) *GetGroupController {
	return &GetGroupController{useCase: useCase}
}

func (ctrl *GetGroupController) Handle(c *gin.Context) {
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

	group, err := ctrl.useCase.Execute(groupID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"owner_id":    group.OwnerID,
		"is_active":   group.IsActive,
		"created_at":  group.CreatedAt,
		"updated_at":  group.UpdatedAt,
	})
}
