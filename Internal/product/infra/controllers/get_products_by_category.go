package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/app"
	"github.com/gin-gonic/gin"
)

type GetProductsByCategoryController struct {
	useCase *app.GetProductsByCategory
}

func NewGetProductsByCategoryController(useCase *app.GetProductsByCategory) *GetProductsByCategoryController {
	return &GetProductsByCategoryController{useCase: useCase}
}

func (ctrl *GetProductsByCategoryController) Handle(c *gin.Context) {
	categoryIDParam := c.Param("categoryId")

	categoryID, err := strconv.ParseInt(categoryIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}

	products, err := ctrl.useCase.Execute(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}

	var response []gin.H
	for _, p := range products {
		response = append(response, gin.H{
			"id":            p.ID,
			"category_id":   p.CategoryID,
			"name":          p.Name,
			"presentation":  p.Presentation,
			"size":          p.Size,
			"default_price": p.DefaultPrice,
			"is_predefined": p.IsPredefined,
			"created_at":    p.CreatedAt,
		})
	}

	if response == nil {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, response)
}
