package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/gin-gonic/gin"
)


type ConfirmParticipantPaymentController struct {
	useCase *app.ConfirmParticipantPaymentUseCase
}

func NewConfirmParticipantPaymentController(useCase *app.ConfirmParticipantPaymentUseCase) *ConfirmParticipantPaymentController {
	return &ConfirmParticipantPaymentController{useCase: useCase}
}

func (c *ConfirmParticipantPaymentController) Handle(ctx *gin.Context) {
	// Parse URL params
	outingIDStr := ctx.Param("outing_id")
	outingID, err := strconv.ParseInt(outingIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid routing id"})
		return
	}

	participantIDStr := ctx.Param("participant_id")
	participantID, err := strconv.ParseInt(participantIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant id"})
		return
	}

	// Get authenticated user (who confirms the payment)
	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	confirmedByUserID, ok := userIDRaw.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Execute use case
	payment, err := c.useCase.Execute(outingID, participantID, confirmedByUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "payment confirmed successfully",
		"data":    payment,
	})
}
