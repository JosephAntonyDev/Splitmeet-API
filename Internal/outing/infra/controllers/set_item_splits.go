package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type SetItemSplitsController struct {
	useCase *app.SetItemSplitsUseCase
}

func NewSetItemSplitsController(useCase *app.SetItemSplitsUseCase) *SetItemSplitsController {
	return &SetItemSplitsController{useCase: useCase}
}

func (c *SetItemSplitsController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	itemIDStr := ctx.Param("itemId")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req app.SetItemSplitsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	splits, err := c.useCase.Execute(itemID, userID.(int64), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, splits)
}
