package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type InviteMemberController struct {
	useCase *app.InviteMember
}

func NewInviteMemberController(useCase *app.InviteMember) *InviteMemberController {
	return &InviteMemberController{useCase: useCase}
}

type InviteMemberRequest struct {
	Username string `json:"username" binding:"required"`
}

func (ctrl *InviteMemberController) Handle(c *gin.Context) {
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

	var req InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El username es requerido"})
		return
	}

	input := app.InviteMemberInput{
		GroupID:   groupID,
		Username:  req.Username,
		InviterID: userID.(int64),
	}

	member, err := ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         member.ID,
		"group_id":   member.GroupID,
		"user_id":    member.UserID,
		"status":     member.Status,
		"invited_by": member.InvitedBy,
		"invited_at": member.InvitedAt,
		"message":    "Invitación enviada exitosamente",
	})
}
