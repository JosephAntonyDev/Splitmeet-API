package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type RemoveParticipantController struct {
	useCase *app.RemoveParticipantUseCase
}

func NewRemoveParticipantController(useCase *app.RemoveParticipantUseCase) *RemoveParticipantController {
	return &RemoveParticipantController{useCase: useCase}
}

func (c *RemoveParticipantController) Handle(ctx *gin.Context) {
	removerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	outingIDStr := ctx.Param("id")
	outingID, err := strconv.ParseInt(outingIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid outing id"})
		return
	}

	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = c.useCase.Execute(outingID, userID, int64(removerID.(int)))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "participant removed successfully"})
}
