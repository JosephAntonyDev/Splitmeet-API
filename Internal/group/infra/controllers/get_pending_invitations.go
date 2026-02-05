package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type GetPendingInvitationsController struct {
	useCase *app.GetPendingInvitations
}

func NewGetPendingInvitationsController(useCase *app.GetPendingInvitations) *GetPendingInvitationsController {
	return &GetPendingInvitationsController{useCase: useCase}
}

func (ctrl *GetPendingInvitationsController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	invitations, err := ctrl.useCase.Execute(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener invitaciones"})
		return
	}

	var response []gin.H
	for _, inv := range invitations {
		response = append(response, gin.H{
			"id":         inv.ID,
			"group_id":   inv.GroupID,
			"user_id":    inv.UserID,
			"status":     inv.Status,
			"invited_by": inv.InvitedBy,
			"invited_at": inv.InvitedAt,
		})
	}

	if response == nil {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, response)
}
