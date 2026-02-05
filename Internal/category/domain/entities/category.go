package entities

import "time"

type Category struct {
	ID        int64
	Name      string
	Icon      string
	IsActive  bool
	CreatedAt time.Time
}
