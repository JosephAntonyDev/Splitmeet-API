package entities

type ItemSplit struct {
	ID            int64
	OutingItemID  int64
	ParticipantID int64
	SplitAmount   float64
	Percentage    *float64
}

// ItemSplitWithUser incluye datos del usuario
type ItemSplitWithUser struct {
	ItemSplit
	UserID   int64
	Username string
	Name     string
}
