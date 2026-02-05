package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/repository"
)

type SearchProducts struct {
	repo repository.ProductRepository
}

func NewSearchProducts(repo repository.ProductRepository) *SearchProducts {
	return &SearchProducts{repo: repo}
}

func (uc *SearchProducts) Execute(query string, categoryID *int64) ([]entities.Product, error) {
	products, err := uc.repo.Search(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar productos: %v", err)
	}
	return products, nil
}
