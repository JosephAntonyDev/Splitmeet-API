package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/category/app"
	"github.com/gin-gonic/gin"
)

type GetAllCategoriesController struct {
	useCase *app.GetAllCategories
}

func NewGetAllCategoriesController(useCase *app.GetAllCategories) *GetAllCategoriesController {
	return &GetAllCategoriesController{useCase: useCase}
}

func (ctrl *GetAllCategoriesController) Handle(c *gin.Context) {
	categories, err := ctrl.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
		return
	}

	var response []gin.H
	for _, cat := range categories {
		response = append(response, gin.H{
			"id":         cat.ID,
			"name":       cat.Name,
			"icon":       cat.Icon,
			"created_at": cat.CreatedAt,
		})
	}

	if response == nil {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, response)
}
