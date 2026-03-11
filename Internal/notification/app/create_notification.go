package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/repository"
)

type CreateNotification struct {
	repo repository.NotificationRepository
}

func NewCreateNotification(repo repository.NotificationRepository) *CreateNotification {
	return &CreateNotification{repo: repo}
}

type CreateNotificationInput struct {
	UserID      int64
	Type        entities.NotificationType
	Title       string
	Message     string
	ReferenceID *int64
	InviterName string
	GroupName   string
	OutingName  string
}

func (uc *CreateNotification) Execute(input CreateNotificationInput) (*entities.Notification, error) {
	notification := &entities.Notification{
		UserID:      input.UserID,
		Type:        input.Type,
		Title:       input.Title,
		Message:     input.Message,
		ReferenceID: input.ReferenceID,
		InviterName: input.InviterName,
		GroupName:   input.GroupName,
		OutingName:  input.OutingName,
		IsRead:      false,
	}

	err := uc.repo.Save(notification)
	if err != nil {
		return nil, fmt.Errorf("error al crear notificación: %v", err)
	}

	return notification, nil
}
