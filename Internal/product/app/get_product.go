package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/repository"
)

type GetProduct struct {
	repo repository.ProductRepository
}

func NewGetProduct(repo repository.ProductRepository) *GetProduct {
	return &GetProduct{repo: repo}
}

func (uc *GetProduct) Execute(id int64) (*entities.Product, error) {
	product, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al buscar producto: %v", err)
	}
	if product == nil {
		return nil, fmt.Errorf("producto no encontrado")
	}
	return product, nil
}
