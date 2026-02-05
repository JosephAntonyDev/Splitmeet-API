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

	now := time.Now()
	payment.Status = entities.PaymentStatusPaid
	payment.PaidAt = &now
	payment.ConfirmedBy = &confirmedByUserID
	payment.UpdatedAt = now

	if err := uc.repo.Update(payment); err != nil {
		return nil, err
	}

	return payment, nil
}
