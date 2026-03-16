package app

import (
	"errors"
	"math"
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
	// 1. Verificar que el participante pertenece a la salida
	exists, err := uc.repo.IsParticipantInOuting(outingID, participantID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("participant not found in this outing")
	}

	// 2. Buscar pago pendiente existente
	payment, err := uc.repo.GetPendingByOutingAndParticipant(outingID, participantID)
	if err != nil && err.Error() != "no pending payment found for this participant in this outing" {
		return nil, err
	}

	// 3. Si no existe pago pendiente, auto-crear uno
	if payment == nil {
		amount, err := uc.calculatePaymentAmount(outingID, participantID)
		if err != nil {
			return nil, err
		}
		if amount <= 0 {
			return nil, errors.New("outing has no amount to pay yet")
		}

		now := time.Now()
		payment = &entities.Payment{
			OutingID:      outingID,
			ParticipantID: participantID,
			Amount:        amount,
			Status:        entities.PaymentStatusPending,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := uc.repo.Create(payment); err != nil {
			return nil, err
		}
	}

	// 4. Validaciones de estado
	if payment.Status == entities.PaymentStatusPaid {
		return nil, errors.New("payment already confirmed")
	}
	if payment.Status == entities.PaymentStatusCancelled {
		return nil, errors.New("payment was cancelled")
	}

	// 5. Confirmar el pago
	now := time.Now()
	payment.Status = entities.PaymentStatusPaid
	payment.PaidAt = &now
	payment.ConfirmedBy = &confirmedByUserID
	payment.UpdatedAt = now

	if err := uc.repo.Update(payment); err != nil {
		return nil, err
	}

	// 6. Verificar si el outing ya está completamente pagado
	outingTotal, err := uc.repo.GetOutingTotalAmount(payment.OutingID)
	if err != nil {
		return payment, nil
	}

	confirmedPayments, err := uc.repo.GetTotalConfirmedPayments(payment.OutingID)
	if err != nil {
		return payment, nil
	}

	// 7. Si el total ya fue alcanzado o superado, cancelar pagos pendientes
	if confirmedPayments >= outingTotal {
		uc.repo.CancelPendingPayments(payment.OutingID)
	}

	return payment, nil
}

// calculatePaymentAmount determina el monto del pago.
// Primero intenta usar amount_owed del participante.
// Si es 0, calcula un split equitativo: total / participantes confirmados.
func (uc *ConfirmParticipantPaymentUseCase) calculatePaymentAmount(outingID, participantID int64) (float64, error) {
	// Intentar obtener el amount_owed individual
	amountOwed, err := uc.repo.GetParticipantAmountOwed(outingID, participantID)
	if err != nil {
		return 0, err
	}
	if amountOwed > 0 {
		return math.Round(amountOwed*100) / 100, nil
	}

	// Fallback: split equitativo
	outingTotal, err := uc.repo.GetOutingTotalAmount(outingID)
	if err != nil {
		return 0, err
	}
	if outingTotal <= 0 {
		return 0, nil
	}

	participantCount, err := uc.repo.GetConfirmedParticipantCount(outingID)
	if err != nil || participantCount == 0 {
		return 0, errors.New("no confirmed participants in outing")
	}

	perPerson := outingTotal / float64(participantCount)
	return math.Round(perPerson*100) / 100, nil
}
