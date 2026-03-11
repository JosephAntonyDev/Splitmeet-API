package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type GetMembersController struct {
	useCase *app.GetMembers
}

func NewGetMembersController(useCase *app.GetMembers) *GetMembersController {
	return &GetMembersController{useCase: useCase}
}

func (ctrl *GetMembersController) Handle(c *gin.Context) {
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

	members, err := ctrl.useCase.Execute(groupID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response []gin.H
	for _, m := range members {
		response = append(response, gin.H{
			"id":           m.ID,
			"user_id":      m.UserID,
			"username":     m.Username,
			"name":         m.Name,
			"email":        m.Email,
			"role":         m.Role,
			"status":       m.Status,
			"invited_by":   m.InvitedBy,
			"invited_at":   m.InvitedAt,
			"responded_at": m.RespondedAt,
		})
	}

	if response == nil {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, response)
}
