package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/repository"
)

type MarkAsRead struct {
	repo repository.NotificationRepository
}

func NewMarkAsRead(repo repository.NotificationRepository) *MarkAsRead {
	return &MarkAsRead{repo: repo}
}

func (uc *MarkAsRead) Execute(notificationID, userID int64) error {
	err := uc.repo.MarkAsRead(notificationID, userID)
	if err != nil {
		return fmt.Errorf("error al marcar como leída: %v", err)
	}
	return nil
}

type MarkAllAsRead struct {
	repo repository.NotificationRepository
}

func NewMarkAllAsRead(repo repository.NotificationRepository) *MarkAllAsRead {
	return &MarkAllAsRead{repo: repo}
}

func (uc *MarkAllAsRead) Execute(userID int64) error {
	err := uc.repo.MarkAllAsRead(userID)
	if err != nil {
		return fmt.Errorf("error al marcar todas como leídas: %v", err)
	}
	return nil
}
