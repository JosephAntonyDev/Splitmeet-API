package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/app"
	"github.com/gin-gonic/gin"
)

type SearchProductsController struct {
	useCase *app.SearchProducts
}

func NewSearchProductsController(useCase *app.SearchProducts) *SearchProductsController {
	return &SearchProductsController{useCase: useCase}
}

func (ctrl *SearchProductsController) Handle(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro de búsqueda 'q' es requerido"})
		return
	}

	var categoryID *int64
	if categoryIDParam := c.Query("category_id"); categoryIDParam != "" {
		id, err := strconv.ParseInt(categoryIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
			return
		}
		categoryID = &id
	}

	products, err := ctrl.useCase.Execute(query, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar productos"})
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
