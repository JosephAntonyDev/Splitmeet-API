package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"

type PaymentRepository interface {
	// Payment CRUD
	Create(payment *entities.Payment) error
	GetByID(id int64) (*entities.Payment, error)
	GetByIDWithDetails(id int64) (*entities.PaymentWithDetails, error)
	Update(payment *entities.Payment) error
	Delete(id int64) error

	// Queries
	GetByOutingID(outingID int64) ([]entities.PaymentWithDetails, error)
	GetByParticipantID(participantID int64) ([]entities.PaymentWithDetails, error)
	GetPendingByOutingID(outingID int64) ([]entities.PaymentWithDetails, error)
	GetPendingByOutingAndParticipant(outingID, participantID int64) (*entities.Payment, error)

	// Summary
	GetSummaryByOutingID(outingID int64) (*entities.PaymentSummary, error)

	// Validations
	GetParticipantIDByOutingAndUser(outingID, userID int64) (int64, error)
	IsParticipantInOuting(outingID, participantID int64) (bool, error)
	GetParticipantAmountOwed(outingID, participantID int64) (float64, error)
	GetConfirmedParticipantCount(outingID int64) (int, error)

	// Outing totals
	GetOutingTotalAmount(outingID int64) (float64, error)
	GetTotalConfirmedPayments(outingID int64) (float64, error)
	CancelPendingPayments(outingID int64) error
}
