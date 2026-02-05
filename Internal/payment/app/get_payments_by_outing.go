package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type GetPaymentsByOutingUseCase struct {
	repo repository.PaymentRepository
}

func NewGetPaymentsByOutingUseCase(repo repository.PaymentRepository) *GetPaymentsByOutingUseCase {
	return &GetPaymentsByOutingUseCase{repo: repo}
}

func (uc *GetPaymentsByOutingUseCase) Execute(outingID int64) ([]entities.PaymentWithDetails, error) {
	return uc.repo.GetByOutingID(outingID)
}
