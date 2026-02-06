package entities

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID            int64
	OutingID      int64
	ParticipantID int64 // Referencia a outing_participants
	Amount        float64
	Status        PaymentStatus
	PaidAt        *time.Time
	ConfirmedBy   *int64 // Usuario que confirmó el pago
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PaymentWithDetails incluye información adicional
type PaymentWithDetails struct {
	Payment
	OutingName          string
	ParticipantUsername string
	ParticipantName     string
	ConfirmedByUsername string
}

// PaymentSummary resumen de pagos para un outing
type PaymentSummary struct {
	OutingID      int64
	TotalAmount   float64
	PaidAmount    float64
	PendingAmount float64
	PaymentsCount int
	PaidCount     int
	PendingCount  int
}
