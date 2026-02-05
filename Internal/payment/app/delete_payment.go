package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type DeletePaymentUseCase struct {
	repo repository.PaymentRepository
}

func NewDeletePaymentUseCase(repo repository.PaymentRepository) *DeletePaymentUseCase {
	return &DeletePaymentUseCase{repo: repo}
}

func (uc *DeletePaymentUseCase) Execute(paymentID int64) error {
	payment, err := uc.repo.GetByID(paymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	// Solo se pueden eliminar pagos pendientes
	if payment.Status == entities.PaymentStatusPaid {
		return errors.New("cannot delete a confirmed payment")
	}

	return uc.repo.Delete(paymentID)
}
