package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/repository"
)

type CreateCustomProduct struct {
	repo repository.ProductRepository
}

func NewCreateCustomProduct(repo repository.ProductRepository) *CreateCustomProduct {
	return &CreateCustomProduct{repo: repo}
}

type CreateCustomProductInput struct {
	CategoryID   *int64
	Name         string
	Presentation string
	Size         string
	DefaultPrice *float64
	CreatedBy    int64
}

func (uc *CreateCustomProduct) Execute(input CreateCustomProductInput) (*entities.Product, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("el nombre del producto es requerido")
	}

	product := &entities.Product{
		CategoryID:   input.CategoryID,
		Name:         input.Name,
		Presentation: input.Presentation,
		Size:         input.Size,
		DefaultPrice: input.DefaultPrice,
		IsPredefined: false,
		CreatedBy:    &input.CreatedBy,
		CreatedAt:    time.Now(),
	}

	err := uc.repo.Save(product)
	if err != nil {
		return nil, fmt.Errorf("error al crear producto: %v", err)
	}

	return product, nil
}
