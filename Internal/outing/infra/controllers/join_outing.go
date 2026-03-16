package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type JoinOutingController struct {
	useCase *app.JoinOutingUseCase
}

func NewJoinOutingController(useCase *app.JoinOutingUseCase) *JoinOutingController {
	return &JoinOutingController{useCase: useCase}
}

func (c *JoinOutingController) Handle(ctx *gin.Context) {
	// Parse outing ID from URL
	outingIDStr := ctx.Param("id")
	outingID, err := strconv.ParseInt(outingIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid outing id"})
		return
	}

	// Get authenticated user ID
	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Execute use case
	participant, err := c.useCase.Execute(outingID, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully joined outing",
		"data":    participant,
	})
}
