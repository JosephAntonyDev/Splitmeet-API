package entities

import "time"

type MemberStatus string

const (
	MemberStatusPending  MemberStatus = "pending"
	MemberStatusAccepted MemberStatus = "accepted"
	MemberStatusRejected MemberStatus = "rejected"
)

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "owner"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

type GroupMember struct {
	ID          int64
	GroupID     int64
	UserID      int64
	Role        MemberRole
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
