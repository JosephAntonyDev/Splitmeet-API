package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type CreatePaymentUseCase struct {
	repo repository.PaymentRepository
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{repo: repo}
}

type CreatePaymentRequest struct {
	OutingID int64   `json:"outing_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Notes    string  `json:"notes"`
}

func (uc *CreatePaymentUseCase) Execute(userID int64, req CreatePaymentRequest) (*entities.Payment, error) {
	// Obtener el participant_id del usuario en este outing
	participantID, err := uc.repo.GetParticipantIDByOutingAndUser(req.OutingID, userID)
	if err != nil {
		return nil, errors.New("user is not a participant in this outing")
	}

	// Obtener el total de la salida
	outingTotal, err := uc.repo.GetOutingTotalAmount(req.OutingID)
	if err != nil {
		return nil, err
	}

	if outingTotal == 0 {
		return nil, errors.New("outing has no items yet, cannot register payment")
	}

	// Obtener la suma de pagos confirmados
	confirmedPayments, err := uc.repo.GetTotalConfirmedPayments(req.OutingID)
	if err != nil {
		return nil, err
	}

	// Calcular el monto restante por pagar
	remainingAmount := outingTotal - confirmedPayments

	if remainingAmount <= 0 {
		return nil, errors.New("outing is already fully paid")
	}

	// Validar que el pago no exceda el monto restante
	if req.Amount > remainingAmount {
		return nil, fmt.Errorf("payment amount (%.2f) exceeds remaining balance (%.2f)", req.Amount, remainingAmount)
	}

	payment := &entities.Payment{
		OutingID:      req.OutingID,
		ParticipantID: participantID,
		Amount:        req.Amount,
		Status:        entities.PaymentStatusPending,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}
