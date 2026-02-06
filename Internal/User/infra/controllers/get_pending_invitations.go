package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
	"github.com/gin-gonic/gin"
)

type GetPendingInvitationsController struct {
	useCase *app.GetPendingInvitations
}

func NewGetPendingInvitationsController(useCase *app.GetPendingInvitations) *GetPendingInvitationsController {
	return &GetPendingInvitationsController{useCase: useCase}
}

func (c *GetPendingInvitationsController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	invitations, err := c.useCase.Execute(userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, invitations)
}
