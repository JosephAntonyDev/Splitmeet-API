package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/gin-gonic/gin"
)

type GetPaymentController struct {
	useCase *app.GetPaymentUseCase
}

func NewGetPaymentController(useCase *app.GetPaymentUseCase) *GetPaymentController {
	return &GetPaymentController{useCase: useCase}
}

func (c *GetPaymentController) Handle(ctx *gin.Context) {
	_, exists := ctx.Get("userID")
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

	payment, err := c.useCase.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
