package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type GetMyGroupsController struct {
	useCase *app.GetMyGroups
}

func NewGetMyGroupsController(useCase *app.GetMyGroups) *GetMyGroupsController {
	return &GetMyGroupsController{useCase: useCase}
}

func (ctrl *GetMyGroupsController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	groups, err := ctrl.useCase.Execute(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener grupos"})
		return
	}

	var response []gin.H
	for _, g := range groups {
		response = append(response, gin.H{
			"id":          g.ID,
			"name":        g.Name,
			"description": g.Description,
			"owner_id":    g.OwnerID,
			"created_at":  g.CreatedAt,
		})
	}

	if response == nil {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, response)
}
