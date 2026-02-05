package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/category/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/category/domain/repository"
)

type GetAllCategories struct {
	repo repository.CategoryRepository
}

func NewGetAllCategories(repo repository.CategoryRepository) *GetAllCategories {
	return &GetAllCategories{repo: repo}
}

func (uc *GetAllCategories) Execute() ([]entities.Category, error) {
	categories, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener categorías: %v", err)
	}
	return categories, nil
}
