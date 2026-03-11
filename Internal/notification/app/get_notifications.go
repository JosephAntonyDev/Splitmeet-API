package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/repository"
)

type GetNotifications struct {
	repo repository.NotificationRepository
}

func NewGetNotifications(repo repository.NotificationRepository) *GetNotifications {
	return &GetNotifications{repo: repo}
}

type GetNotificationsInput struct {
	UserID int64
	Limit  int
	Offset int
}

func (uc *GetNotifications) Execute(input GetNotificationsInput) ([]entities.Notification, int, error) {
	notifications, total, err := uc.repo.GetByUserID(input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener notificaciones: %v", err)
	}
	return notifications, total, nil
}
