package entities

import "time"

type Group struct {
	ID          int64
	Name        string
	Description string
	OwnerID     int64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GroupWithDetails struct {
	Group
	OwnerUsername string
	MemberCount   int
}
