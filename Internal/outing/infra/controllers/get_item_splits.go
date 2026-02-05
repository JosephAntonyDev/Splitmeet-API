package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type GetItemSplitsController struct {
	useCase *app.GetItemSplitsUseCase
}

func NewGetItemSplitsController(useCase *app.GetItemSplitsUseCase) *GetItemSplitsController {
	return &GetItemSplitsController{useCase: useCase}
}

func (c *GetItemSplitsController) Handle(ctx *gin.Context) {
	_, exists := ctx.Get("userID")
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

	splits, err := c.useCase.Execute(itemID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, splits)
}
