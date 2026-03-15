package core

import (
	"context"
	"fmt"
	"log"
	"time"
)

// NotificationService provides a shared way for any module to create notifications and push SSE events
type NotificationService struct {
	db         *Conn_PostgreSQL
	hub        *SSEHub
	pushSender PushSender
	tokenStore DeviceTokenStore
}

type PushRequest struct {
	Token string
	Title string
	Body  string
	Data  map[string]string
}

type PushSender interface {
	SendAndroidPush(ctx context.Context, req PushRequest) (bool, error)
}

type DeviceTokenStore interface {
	GetActiveDeviceTokensByUserID(userID int64) ([]string, error)
	DeactivateDeviceToken(token string) error
}

func NewNotificationService(db *Conn_PostgreSQL, hub *SSEHub, pushSender PushSender, tokenStore DeviceTokenStore) *NotificationService {
	return &NotificationService{db: db, hub: hub, pushSender: pushSender, tokenStore: tokenStore}
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
	if s == nil {
		return
	}

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

	if s.pushSender == nil || s.tokenStore == nil {
		log.Printf("⚠️ FCM push saltado para usuario %d: pushSender=%v, tokenStore=%v", payload.UserID, s.pushSender != nil, s.tokenStore != nil)
		return
	}

	tokens, err := s.tokenStore.GetActiveDeviceTokensByUserID(payload.UserID)
	if err != nil {
		log.Printf("Error al obtener tokens FCM de usuario %d: %v", payload.UserID, err)
		return
	}

	if len(tokens) == 0 {
		log.Printf("⚠️ Usuario %d no tiene tokens FCM activos — no se envía push", payload.UserID)
		return
	}

	log.Printf("📱 Enviando FCM push a usuario %d (%d tokens)", payload.UserID, len(tokens))

	for _, token := range tokens {
		invalidToken, pushErr := s.pushSender.SendAndroidPush(context.Background(), PushRequest{
			Token: token,
			Title: payload.Title,
			Body:  payload.Message,
			Data: map[string]string{
				"type":    payload.Type,
				"user_id": fmt.Sprintf("%d", payload.UserID),
			},
		})

		if pushErr != nil {
			log.Printf("❌ Error enviando push a usuario %d: %v", payload.UserID, pushErr)
			if invalidToken {
				if deactivateErr := s.tokenStore.DeactivateDeviceToken(token); deactivateErr != nil {
					log.Printf("Error desactivando token inválido: %v", deactivateErr)
				}
			}
		} else {
			log.Printf("✅ Push FCM enviado exitosamente a usuario %d", payload.UserID)
		}
	}
}
