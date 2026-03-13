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

	// Query que une con outing_participants y group_members para obtener el response_status
	query := `
		SELECT n.id, n.user_id, n.type, n.title, n.message, n.reference_id, 
			   COALESCE(n.inviter_name, ''), COALESCE(n.group_name, ''), COALESCE(n.outing_name, ''),
			   n.is_read, n.created_at,
			   CASE 
			       WHEN n.type = 'outing_invitation' THEN 
			           COALESCE((SELECT 
			               CASE op.status 
			                   WHEN 'confirmed' THEN 'accepted'
			                   WHEN 'declined' THEN 'rejected'
			                   ELSE 'pending'
			               END
			           FROM outing_participants op 
			           WHERE op.outing_id = n.reference_id AND op.user_id = n.user_id), 'pending')
			       WHEN n.type = 'group_invitation' THEN 
			           COALESCE((SELECT 
			               CASE gm.status 
			                   WHEN 'accepted' THEN 'accepted'
			                   WHEN 'rejected' THEN 'rejected'
			                   ELSE 'pending'
			               END
			           FROM group_members gm 
			           WHERE gm.group_id = n.reference_id AND gm.user_id = n.user_id), 'pending')
			       ELSE 'pending'
			   END as response_status
		FROM notifications n
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
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
		var responseStatus string

		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &refID,
			&n.InviterName, &n.GroupName, &n.OutingName,
			&n.IsRead, &n.CreatedAt, &responseStatus,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al escanear notificación: %v", err)
		}
		if refID.Valid {
			n.ReferenceID = &refID.Int64
		}
		n.ResponseStatus = entities.ResponseStatus(responseStatus)
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

func (r *NotificationPostgreSQLRepository) UpsertDeviceToken(userID int64, token, platform string) error {
	query := `
		INSERT INTO user_device_tokens (user_id, token, platform, is_active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (token)
		DO UPDATE SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform, is_active = true, updated_at = CURRENT_TIMESTAMP`

	_, err := r.conn.DB.Exec(query, userID, token, platform)
	if err != nil {
		return fmt.Errorf("error al guardar token de dispositivo: %v", err)
	}

	return nil
}

func (r *NotificationPostgreSQLRepository) GetActiveDeviceTokensByUserID(userID int64) ([]string, error) {
	rows, err := r.conn.DB.Query(`
		SELECT token
		FROM user_device_tokens
		WHERE user_id = $1 AND is_active = true`, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener tokens activos: %v", err)
	}
	defer rows.Close()

	tokens := make([]string, 0)
	for rows.Next() {
		var token string
		if scanErr := rows.Scan(&token); scanErr != nil {
			return nil, fmt.Errorf("error al leer token: %v", scanErr)
		}
		tokens = append(tokens, token)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error al iterar tokens: %v", rowsErr)
	}

	return tokens, nil
}

func (r *NotificationPostgreSQLRepository) DeactivateDeviceToken(token string) error {
	_, err := r.conn.DB.Exec(`
		UPDATE user_device_tokens
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE token = $1`, token)
	if err != nil {
		return fmt.Errorf("error al desactivar token: %v", err)
	}

	return nil
}
