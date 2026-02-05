package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/gin-gonic/gin"
)

type GetParticipantsController struct {
	useCase *app.GetParticipantsUseCase
}

func NewGetParticipantsController(useCase *app.GetParticipantsUseCase) *GetParticipantsController {
	return &GetParticipantsController{useCase: useCase}
}

func (c *GetParticipantsController) Handle(ctx *gin.Context) {
	outingIDStr := ctx.Param("id")
	outingID, err := strconv.ParseInt(outingIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid outing id"})
		return
	}

	participants, err := c.useCase.Execute(outingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if participants == nil {
		participants = []entities.OutingParticipantWithUser{}
	}

	ctx.JSON(http.StatusOK, participants)
}
