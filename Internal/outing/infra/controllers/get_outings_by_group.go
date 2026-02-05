package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/gin-gonic/gin"
)

type GetOutingsByGroupController struct {
	useCase *app.GetOutingsByGroupUseCase
}

func NewGetOutingsByGroupController(useCase *app.GetOutingsByGroupUseCase) *GetOutingsByGroupController {
	return &GetOutingsByGroupController{useCase: useCase}
}

func (c *GetOutingsByGroupController) Handle(ctx *gin.Context) {
	groupIDStr := ctx.Param("groupId")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	outings, err := c.useCase.Execute(groupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if outings == nil {
		outings = []entities.OutingWithDetails{}
	}

	ctx.JSON(http.StatusOK, outings)
}
