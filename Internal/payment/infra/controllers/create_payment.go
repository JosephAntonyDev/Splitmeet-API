package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/gin-gonic/gin"
)

type CreatePaymentController struct {
	useCase *app.CreatePaymentUseCase
}

func NewCreatePaymentController(useCase *app.CreatePaymentUseCase) *CreatePaymentController {
	return &CreatePaymentController{useCase: useCase}
}

func (c *CreatePaymentController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req app.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := c.useCase.Execute(userID.(int64), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, payment)
}
