package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/category/domain/entities"

type CategoryRepository interface {
	GetAll() ([]entities.Category, error)
	GetByID(id int64) (*entities.Category, error)
}
