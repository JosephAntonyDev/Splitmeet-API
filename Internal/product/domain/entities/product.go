package entities

import "time"

type Product struct {
	ID           int64
	CategoryID   *int64
	Name         string
	Presentation string
	Size         string
	DefaultPrice *float64
	IsPredefined bool
	CreatedBy    *int64
	CreatedAt    time.Time
}
