package app

import (
	"errors"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type ConfirmParticipantPaymentUseCase struct {
	repo repository.PaymentRepository
}

func NewConfirmParticipantPaymentUseCase(repo repository.PaymentRepository) *ConfirmParticipantPaymentUseCase {
	return &ConfirmParticipantPaymentUseCase{repo: repo}
}

func (uc *ConfirmParticipantPaymentUseCase) Execute(outingID, participantID, confirmedByUserID int64) (*entities.Payment, error) {
	// 1. Buscar el pago pendiente del participante en la salida
	payment, err := uc.repo.GetPendingByOutingAndParticipant(outingID, participantID)
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

	// 2. Actualizar el estado del pago
	if err := uc.repo.Update(payment); err != nil {
		return nil, err
	}

	// 3. Verificar si el outing ya está completamente pagado
	outingTotal, err := uc.repo.GetOutingTotalAmount(payment.OutingID)
	if err != nil {
		return payment, nil // El pago se confirmó, pero no pudimos verificar el total
	}

	confirmedPayments, err := uc.repo.GetTotalConfirmedPayments(payment.OutingID)
	if err != nil {
		return payment, nil
	}

	// 4. Si el total ya fue alcanzado o superado, cancelar pagos pendientes
	if confirmedPayments >= outingTotal {
		uc.repo.CancelPendingPayments(payment.OutingID)
	}

	return payment, nil
}
