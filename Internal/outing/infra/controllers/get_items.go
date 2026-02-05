package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/gin-gonic/gin"
)

type GetItemsController struct {
	useCase *app.GetItemsUseCase
}

func NewGetItemsController(useCase *app.GetItemsUseCase) *GetItemsController {
	return &GetItemsController{useCase: useCase}
}

func (c *GetItemsController) Handle(ctx *gin.Context) {
	outingIDStr := ctx.Param("id")
	outingID, err := strconv.ParseInt(outingIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid outing id"})
		return
	}

	items, err := c.useCase.Execute(outingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		items = []entities.OutingItemWithProduct{}
	}

	ctx.JSON(http.StatusOK, items)
}
