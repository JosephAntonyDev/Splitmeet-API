package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"

type UserRepository interface {
	Save(user *entities.User) error
	GetByID(id int64) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	GetUsersByIDs(ids []int64) ([]entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	SearchByUsername(query string, limit int) ([]entities.User, error)
	Update(user *entities.User) error
	Delete(id int64) error
}
