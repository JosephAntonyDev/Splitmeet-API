package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type GetProfile struct {
	repo repository.UserRepository
}

func NewGetProfile(repo repository.UserRepository) *GetProfile {
	return &GetProfile{repo: repo}
}

func (uc *GetProfile) Execute(id int64) (*entities.User, error) {
	user, err := uc.repo.GetByID(id)
	
	if err != nil {
		return nil, fmt.Errorf("error del sistema al buscar usuario: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	user.Password = "" 

	return user, nil
}