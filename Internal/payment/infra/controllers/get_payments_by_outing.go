package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/gin-gonic/gin"
)

type GetPaymentsByOutingController struct {
	useCase *app.GetPaymentsByOutingUseCase
}

func NewGetPaymentsByOutingController(useCase *app.GetPaymentsByOutingUseCase) *GetPaymentsByOutingController {
	return &GetPaymentsByOutingController{useCase: useCase}
}

func (c *GetPaymentsByOutingController) Handle(ctx *gin.Context) {
	_, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	outingIDStr := ctx.Param("outingId")
	outingID, err := strconv.ParseInt(outingIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid outing id"})
		return
	}

	payments, err := c.useCase.Execute(outingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payments)
}
