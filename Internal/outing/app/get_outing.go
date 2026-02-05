package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type GetOutingUseCase struct {
	repo repository.OutingRepository
}

func NewGetOutingUseCase(repo repository.OutingRepository) *GetOutingUseCase {
	return &GetOutingUseCase{repo: repo}
}

func (uc *GetOutingUseCase) Execute(outingID int64) (*entities.OutingWithDetails, error) {
	return uc.repo.GetByIDWithDetails(outingID)
}
