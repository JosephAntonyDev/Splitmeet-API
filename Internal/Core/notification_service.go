package core

import (
	"fmt"
	"time"
)

// NotificationService provides a shared way for any module to create notifications and push SSE events
type NotificationService struct {
	db  *Conn_PostgreSQL
	hub *SSEHub
}

func NewNotificationService(db *Conn_PostgreSQL, hub *SSEHub) *NotificationService {
	return &NotificationService{db: db, hub: hub}
}

type NotificationPayload struct {
	UserID      int64  `json:"user_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	ReferenceID *int64 `json:"reference_id,omitempty"`
	InviterName string `json:"inviter_name,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	OutingName  string `json:"outing_name,omitempty"`
}

func (s *NotificationService) Send(payload NotificationPayload) {
	var id int64
	var createdAt time.Time

	err := s.db.DB.QueryRow(`
		INSERT INTO notifications (user_id, type, title, message, reference_id, inviter_name, group_name, outing_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`,
		payload.UserID, payload.Type, payload.Title, payload.Message,
		payload.ReferenceID, payload.InviterName, payload.GroupName, payload.OutingName,
	).Scan(&id, &createdAt)

	if err != nil {
		fmt.Printf("Error al guardar notificación: %v\n", err)
		return
	}

	// Push SSE event to the user
	s.hub.SendToUser(payload.UserID, "notification", map[string]interface{}{
		"id":              id,
		"type":            payload.Type,
		"title":           payload.Title,
		"message":         payload.Message,
		"reference_id":    payload.ReferenceID,
		"inviter_name":    payload.InviterName,
		"group_name":      payload.GroupName,
		"outing_name":     payload.OutingName,
		"is_read":         false,
		"response_status": "pending",
		"created_at":      createdAt,
	})
}
