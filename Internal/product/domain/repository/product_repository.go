package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/entities"

type ProductRepository interface {
	GetByID(id int64) (*entities.Product, error)
	GetByCategory(categoryID int64) ([]entities.Product, error)
	Search(query string, categoryID *int64) ([]entities.Product, error)
	Save(product *entities.Product) error
}
