package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/gin-gonic/gin"
)

type DeletePaymentController struct {
	useCase *app.DeletePaymentUseCase
}

func NewDeletePaymentController(useCase *app.DeletePaymentUseCase) *DeletePaymentController {
	return &DeletePaymentController{useCase: useCase}
}

func (c *DeletePaymentController) Handle(ctx *gin.Context) {
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

	err = c.useCase.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "payment deleted successfully"})
}
