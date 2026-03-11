package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
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

	params := core.GetPaginationParams(c)

	input := app.GetMyGroupsInput{
		UserID: userID.(int64),
		Limit:  params.Limit,
		Offset: params.Offset,
		Search: params.Search,
	}

	groups, total, err := ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener grupos"})
		return
	}

	var data []gin.H
	for _, g := range groups {
		data = append(data, gin.H{
			"id":             g.ID,
			"name":           g.Name,
			"description":    g.Description,
			"owner_id":       g.OwnerID,
			"owner_username": g.OwnerUsername,
			"member_count":   g.MemberCount,
			"created_at":     g.CreatedAt,
		})
	}

	if data == nil {
		data = []gin.H{}
	}

	c.JSON(http.StatusOK, core.NewPaginatedResponse(data, params.Page, params.Limit, total))
}
