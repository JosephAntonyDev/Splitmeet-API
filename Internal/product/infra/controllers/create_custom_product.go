package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/app"
	"github.com/gin-gonic/gin"
)

type CreateCustomProductController struct {
	useCase *app.CreateCustomProduct
}

func NewCreateCustomProductController(useCase *app.CreateCustomProduct) *CreateCustomProductController {
	return &CreateCustomProductController{useCase: useCase}
}

type CreateProductRequest struct {
	CategoryID   *int64   `json:"category_id"`
	Name         string   `json:"name" binding:"required"`
	Presentation string   `json:"presentation"`
	Size         string   `json:"size"`
	DefaultPrice *float64 `json:"default_price"`
}

func (ctrl *CreateCustomProductController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: el nombre es requerido"})
		return
	}

	input := app.CreateCustomProductInput{
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Presentation: req.Presentation,
		Size:         req.Size,
		DefaultPrice: req.DefaultPrice,
		CreatedBy:    userID.(int64),
	}

	product, err := ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
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
