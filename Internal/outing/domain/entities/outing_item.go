package entities

import "time"

type OutingItem struct {
	ID                 int64
	OutingID           int64
	ProductID          *int64
	CustomName         string
	CustomPresentation string
	Quantity           int
	UnitPrice          float64
	Subtotal           float64 // quantity * unit_price
	IsShared           bool
	CreatedAt          time.Time
}

// OutingItemWithProduct incluye el nombre del producto si existe
type OutingItemWithProduct struct {
	OutingItem
	ProductName         string
	ProductPresentation string
	ProductSize         string
}
