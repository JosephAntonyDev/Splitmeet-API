package app

import (
	"fmt"
	"strings"

	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/repository"
)

type RegisterDeviceToken struct {
	repo repository.NotificationRepository
}

func NewRegisterDeviceToken(repo repository.NotificationRepository) *RegisterDeviceToken {
	return &RegisterDeviceToken{repo: repo}
}

type RegisterDeviceTokenInput struct {
	UserID   int64
	Token    string
	Platform string
}

func (uc *RegisterDeviceToken) Execute(input RegisterDeviceTokenInput) error {
	token := strings.TrimSpace(input.Token)
	if token == "" {
		return fmt.Errorf("el token es obligatorio")
	}

	platform := strings.ToLower(strings.TrimSpace(input.Platform))
	if platform == "" {
		platform = "android"
	}

	if platform != "android" {
		return fmt.Errorf("solo se soporta plataforma android")
	}

	if err := uc.repo.UpsertDeviceToken(input.UserID, token, platform); err != nil {
		return err
	}

	return nil
}
