package entities

import "time"

type ParticipantStatus string

const (
	ParticipantStatusPending   ParticipantStatus = "pending"
	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
	ParticipantStatusDeclined  ParticipantStatus = "declined"
)

type OutingParticipant struct {
	ID           int64
	OutingID     int64
	UserID       int64
	InvitedBy    *int64
	Status       ParticipantStatus
	AmountOwed   float64
	CustomAmount *float64
	JoinedAt     time.Time
}

// OutingParticipantWithUser incluye datos del usuario
type OutingParticipantWithUser struct {
	OutingParticipant
	Username string
	Name     string
	Email    string
}
