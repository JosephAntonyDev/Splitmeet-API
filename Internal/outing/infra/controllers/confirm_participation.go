package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type ConfirmParticipationController struct {
	useCase *app.ConfirmParticipationUseCase
}

func NewConfirmParticipationController(useCase *app.ConfirmParticipationUseCase) *ConfirmParticipationController {
	return &ConfirmParticipationController{useCase: useCase}
}

type ConfirmParticipationRequest struct {
	Accept bool `json:"accept"`
}

func (c *ConfirmParticipationController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
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

	var req ConfirmParticipationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.useCase.Execute(outingID, userID.(int64), req.Accept)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "participation updated successfully"})
}
