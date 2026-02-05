package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/category/app"
	"github.com/gin-gonic/gin"
)

type GetCategoryController struct {
	useCase *app.GetCategory
}

func NewGetCategoryController(useCase *app.GetCategory) *GetCategoryController {
	return &GetCategoryController{useCase: useCase}
}

func (ctrl *GetCategoryController) Handle(c *gin.Context) {
	idParam := c.Param("id")

	categoryID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}

	category, err := ctrl.useCase.Execute(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         category.ID,
		"name":       category.Name,
		"icon":       category.Icon,
		"created_at": category.CreatedAt,
	})
}
