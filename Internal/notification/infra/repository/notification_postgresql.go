package repository

import (
	"database/sql"
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/domain/entities"
)

type NotificationPostgreSQLRepository struct {
	conn *core.Conn_PostgreSQL
}

func NewNotificationPostgreSQLRepository(conn *core.Conn_PostgreSQL) *NotificationPostgreSQLRepository {
	return &NotificationPostgreSQLRepository{conn: conn}
}

func (r *NotificationPostgreSQLRepository) Save(notification *entities.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message, reference_id, inviter_name, group_name, outing_name, is_read)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	err := r.conn.DB.QueryRow(
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.ReferenceID,
		notification.InviterName,
		notification.GroupName,
		notification.OutingName,
		notification.IsRead,
	).Scan(&notification.ID, &notification.CreatedAt)

	if err != nil {
		return fmt.Errorf("error al insertar notificación: %v", err)
	}
	return nil
}

func (r *NotificationPostgreSQLRepository) GetByUserID(userID int64, limit, offset int) ([]entities.Notification, int, error) {
	var total int
	err := r.conn.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error al contar notificaciones: %v", err)
	}

	query := `
		SELECT id, user_id, type, title, message, reference_id, 
			   COALESCE(inviter_name, ''), COALESCE(group_name, ''), COALESCE(outing_name, ''),
			   is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.conn.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener notificaciones: %v", err)
	}
	defer rows.Close()

	var notifications []entities.Notification
	for rows.Next() {
		var n entities.Notification
		var refID sql.NullInt64

		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &refID,
			&n.InviterName, &n.GroupName, &n.OutingName,
			&n.IsRead, &n.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al escanear notificación: %v", err)
		}
		if refID.Valid {
			n.ReferenceID = &refID.Int64
		}
		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (r *NotificationPostgreSQLRepository) MarkAsRead(notificationID, userID int64) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	_, err := r.conn.DB.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("error al marcar como leída: %v", err)
	}
	return nil
}

func (r *NotificationPostgreSQLRepository) MarkAllAsRead(userID int64) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.conn.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error al marcar todas como leídas: %v", err)
	}
	return nil
}

func (r *NotificationPostgreSQLRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	err := r.conn.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar no leídas: %v", err)
	}
	return count, nil
}
