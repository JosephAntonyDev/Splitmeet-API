package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type CalculateSplitsController struct {
	useCase *app.CalculateSplitsUseCase
}

func NewCalculateSplitsController(useCase *app.CalculateSplitsUseCase) *CalculateSplitsController {
	return &CalculateSplitsController{useCase: useCase}
}

func (c *CalculateSplitsController) Handle(ctx *gin.Context) {
	_, exists := ctx.Get("userID")
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

	summary, err := c.useCase.Execute(outingID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}
