package app

import (
	"fmt"
	"strings"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/ports"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type CreateUser struct {
	repo   repository.UserRepository
	bcrypt ports.IBcryptService
}

func NewCreateUser(repo repository.UserRepository, bcrypt ports.IBcryptService) *CreateUser {
	return &CreateUser{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

func (cu *CreateUser) Execute(user *entities.User) error {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Name = strings.TrimSpace(user.Name)
	user.Username = strings.TrimSpace(strings.ToLower(user.Username))

	existingUser, _ := cu.repo.GetByEmail(user.Email)
	if existingUser != nil {
		return fmt.Errorf("el usuario con email %s ya existe", user.Email)
	}

	existingByUsername, _ := cu.repo.GetByUsername(user.Username)
	if existingByUsername != nil {
		return fmt.Errorf("el username %s ya está en uso", user.Username)
	}

	hashedPass, err := cu.bcrypt.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error al encriptar contraseña: %v", err)
	}
	user.Password = hashedPass

	if err := cu.repo.Save(user); err != nil {
		return fmt.Errorf("error al guardar usuario: %v", err)
	}

	return nil

}
