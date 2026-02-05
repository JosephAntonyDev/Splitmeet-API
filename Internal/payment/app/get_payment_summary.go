package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/repository"
)

type GetPaymentSummaryUseCase struct {
	repo repository.PaymentRepository
}

func NewGetPaymentSummaryUseCase(repo repository.PaymentRepository) *GetPaymentSummaryUseCase {
	return &GetPaymentSummaryUseCase{repo: repo}
}

func (uc *GetPaymentSummaryUseCase) Execute(outingID int64) (*entities.PaymentSummary, error) {
	return uc.repo.GetSummaryByOutingID(outingID)
}
