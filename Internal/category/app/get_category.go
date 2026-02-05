package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/category/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/category/domain/repository"
)

type GetCategory struct {
	repo repository.CategoryRepository
}

func NewGetCategory(repo repository.CategoryRepository) *GetCategory {
	return &GetCategory{repo: repo}
}

func (uc *GetCategory) Execute(id int64) (*entities.Category, error) {
	category, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al buscar categoría: %v", err)
	}
	if category == nil {
		return nil, fmt.Errorf("categoría no encontrada")
	}
	return category, nil
}
