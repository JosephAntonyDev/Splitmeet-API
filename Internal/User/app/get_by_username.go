package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type GetByUsername struct {
	repo repository.UserRepository
}

func NewGetByUsername(repo repository.UserRepository) *GetByUsername {
	return &GetByUsername{repo: repo}
}

func (uc *GetByUsername) Execute(username string) (*entities.User, error) {
	user, err := uc.repo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	user.Password = ""

	return user, nil
}
