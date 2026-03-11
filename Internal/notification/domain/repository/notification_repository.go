package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/entities"

type NotificationRepository interface {
	Save(notification *entities.Notification) error
	GetByUserID(userID int64, limit, offset int) ([]entities.Notification, int, error)
	MarkAsRead(notificationID, userID int64) error
	MarkAllAsRead(userID int64) error
	GetUnreadCount(userID int64) (int, error)
}
