package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/repository"
)

type GetProductsByCategory struct {
	repo repository.ProductRepository
}

func NewGetProductsByCategory(repo repository.ProductRepository) *GetProductsByCategory {
	return &GetProductsByCategory{repo: repo}
}

func (uc *GetProductsByCategory) Execute(categoryID int64) ([]entities.Product, error) {
	products, err := uc.repo.GetByCategory(categoryID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %v", err)
	}
	return products, nil
}
