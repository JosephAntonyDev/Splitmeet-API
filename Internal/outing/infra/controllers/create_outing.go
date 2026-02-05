package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/gin-gonic/gin"
)

type CreateOutingController struct {
	useCase *app.CreateOutingUseCase
}

func NewCreateOutingController(useCase *app.CreateOutingUseCase) *CreateOutingController {
	return &CreateOutingController{useCase: useCase}
}

func (c *CreateOutingController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req app.CreateOutingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outing, err := c.useCase.Execute(req, userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, outing)
}
