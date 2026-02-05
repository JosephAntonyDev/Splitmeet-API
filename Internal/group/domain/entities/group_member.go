package entities

import "time"

type MemberStatus string

const (
	MemberStatusPending  MemberStatus = "pending"
	MemberStatusAccepted MemberStatus = "accepted"
	MemberStatusRejected MemberStatus = "rejected"
)

type GroupMember struct {
	ID          int64
	GroupID     int64
	UserID      int64
	Status      MemberStatus
	InvitedBy   *int64
	InvitedAt   time.Time
	RespondedAt *time.Time
}

// GroupMemberWithUser incluye datos del usuario para respuestas
type GroupMemberWithUser struct {
	GroupMember
	Username string
	Name     string
	Email    string
}
