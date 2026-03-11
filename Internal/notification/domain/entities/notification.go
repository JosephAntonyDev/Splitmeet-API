package entities

import "time"

type NotificationType string

const (
	NotificationGroupInvitation    NotificationType = "group_invitation"
	NotificationOutingInvitation   NotificationType = "outing_invitation"
	NotificationInvitationAccepted NotificationType = "invitation_accepted"
	NotificationInvitationRejected NotificationType = "invitation_rejected"
)

// ResponseStatus indica el estado de respuesta de una invitación
// "pending" = no respondida, "accepted" = aceptada, "rejected" = rechazada
type ResponseStatus string

const (
	ResponsePending  ResponseStatus = "pending"
	ResponseAccepted ResponseStatus = "accepted"
	ResponseRejected ResponseStatus = "rejected"
)

type Notification struct {
	ID             int64
	UserID         int64
	Type           NotificationType
	Title          string
	Message        string
	ReferenceID    *int64
	InviterName    string
	GroupName      string
	OutingName     string
	IsRead         bool
	ResponseStatus ResponseStatus // Estado de respuesta para invitaciones
	CreatedAt      time.Time
}
