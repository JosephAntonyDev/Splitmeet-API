package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/app"
	"github.com/gin-gonic/gin"
)

type GetProductController struct {
	useCase *app.GetProduct
}

func NewGetProductController(useCase *app.GetProduct) *GetProductController {
	return &GetProductController{useCase: useCase}
}

func (ctrl *GetProductController) Handle(c *gin.Context) {
	idParam := c.Param("id")

	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	product, err := ctrl.useCase.Execute(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            product.ID,
		"category_id":   product.CategoryID,
		"name":          product.Name,
		"presentation":  product.Presentation,
		"size":          product.Size,
		"default_price": product.DefaultPrice,
		"is_predefined": product.IsPredefined,
		"created_by":    product.CreatedBy,
		"created_at":    product.CreatedAt,
	})
}
