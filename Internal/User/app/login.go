package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/ports"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type LoginUser struct {
	repo         repository.UserRepository
	bcrypt       ports.IBcryptService
	tokenManager ports.TokenManager
}

func NewLoginUser(repo repository.UserRepository, bcrypt ports.IBcryptService, tokenManager ports.TokenManager) *LoginUser {
	return &LoginUser{
		repo:         repo,
		bcrypt:       bcrypt,
		tokenManager: tokenManager,
	}
}

func (uc *LoginUser) Execute(email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("error del sistema al buscar usuario: %v", err)
	}
	if user == nil {
		return "", fmt.Errorf("credenciales inválidas")
	}

	match := uc.bcrypt.ComparePasswords(user.Password, password)
	if !match {
		return "", fmt.Errorf("credenciales inválidas")
	}

	token, err := uc.tokenManager.GenerateToken(int(user.ID), user.Email, user.Name) 
	if err != nil {
		return "", fmt.Errorf("error al generar token: %v", err)
	}

	return token, nil
}