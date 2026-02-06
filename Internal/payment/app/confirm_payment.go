package app

import (
	"errors"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type ConfirmPaymentUseCase struct {
	repo repository.PaymentRepository
}

func NewConfirmPaymentUseCase(repo repository.PaymentRepository) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{repo: repo}
}

func (uc *ConfirmPaymentUseCase) Execute(paymentID, confirmedByUserID int64) (*entities.Payment, error) {
	payment, err := uc.repo.GetByID(paymentID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}

	if payment.Status == entities.PaymentStatusPaid {
		return nil, errors.New("payment already confirmed")
	}

	if payment.Status == entities.PaymentStatusCancelled {
		return nil, errors.New("payment was cancelled")
	}

	now := time.Now()
	payment.Status = entities.PaymentStatusPaid
	payment.PaidAt = &now
	payment.ConfirmedBy = &confirmedByUserID
	payment.UpdatedAt = now

	if err := uc.repo.Update(payment); err != nil {
		return nil, err
	}

	// Verificar si el outing ya está completamente pagado
	outingTotal, err := uc.repo.GetOutingTotalAmount(payment.OutingID)
	if err != nil {
		return payment, nil // El pago se confirmó, pero no pudimos verificar el total
	}

	confirmedPayments, err := uc.repo.GetTotalConfirmedPayments(payment.OutingID)
	if err != nil {
		return payment, nil
	}

	// Si el total ya fue alcanzado o superado, cancelar pagos pendientes
	if confirmedPayments >= outingTotal {
		uc.repo.CancelPendingPayments(payment.OutingID)
	}

	return payment, nil
}
