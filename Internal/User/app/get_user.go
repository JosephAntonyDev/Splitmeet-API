package app

import (
	"fmt"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type GetUser struct {
	repo repository.UserRepository
}

func NewGetUser(repo repository.UserRepository) *GetUser {
	return &GetUser{repo: repo}
}

func (uc *GetUser) Execute(id int64) (*entities.User, error) {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	user.Password = "" 
	
	return user, nil
}