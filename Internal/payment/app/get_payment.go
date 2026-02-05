package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type GetPaymentUseCase struct {
	repo repository.PaymentRepository
}

func NewGetPaymentUseCase(repo repository.PaymentRepository) *GetPaymentUseCase {
	return &GetPaymentUseCase{repo: repo}
}

func (uc *GetPaymentUseCase) Execute(paymentID int64) (*entities.PaymentWithDetails, error) {
	return uc.repo.GetByIDWithDetails(paymentID)
}
