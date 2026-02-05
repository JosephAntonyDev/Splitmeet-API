package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/gin-gonic/gin"
)

type ConfirmPaymentController struct {
	useCase *app.ConfirmPaymentUseCase
}

func NewConfirmPaymentController(useCase *app.ConfirmPaymentUseCase) *ConfirmPaymentController {
	return &ConfirmPaymentController{useCase: useCase}
}

func (c *ConfirmPaymentController) Handle(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	payment, err := c.useCase.Execute(id, userID.(int64))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
