package entities

import "time"

type SplitType string

const (
	SplitTypeEqual          SplitType = "equal"           // División equitativa
	SplitTypeCustomFixed    SplitType = "custom_fixed"    // Una persona paga cantidad fija, resto equitativo
	SplitTypePerConsumption SplitType = "per_consumption" // Cada quien lo que pidió
	SplitTypeSinglePayer    SplitType = "single_payer"    // Una persona paga todo
)

type OutingStatus string

const (
	OutingStatusActive    OutingStatus = "active"
	OutingStatusCompleted OutingStatus = "completed"
	OutingStatusCancelled OutingStatus = "cancelled"
)

type Outing struct {
	ID          int64
	Name        string
	Description string
	CategoryID  *int64
	GroupID     *int64
	CreatorID   int64
	OutingDate  time.Time
	SplitType   SplitType
	TotalAmount float64
	Status      OutingStatus
	IsEditable  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OutingWithDetails incluye información adicional para respuestas
type OutingWithDetails struct {
	Outing
	CategoryName     string
	GroupName        string
	CreatorUsername  string
	ParticipantCount int
	PaidCount        int
}
