package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/gin-gonic/gin"
)

type GetMyOutingsController struct {
	useCase *app.GetMyOutingsUseCase
}

func NewGetMyOutingsController(useCase *app.GetMyOutingsUseCase) *GetMyOutingsController {
	return &GetMyOutingsController{useCase: useCase}
}

func (c *GetMyOutingsController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	outings, err := c.useCase.Execute(userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if outings == nil {
		outings = []entities.OutingWithDetails{}
	}

	ctx.JSON(http.StatusOK, outings)
}
