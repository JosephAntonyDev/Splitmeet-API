package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type GetOutingsByGroupUseCase struct {
	repo repository.OutingRepository
}

func NewGetOutingsByGroupUseCase(repo repository.OutingRepository) *GetOutingsByGroupUseCase {
	return &GetOutingsByGroupUseCase{repo: repo}
}

func (uc *GetOutingsByGroupUseCase) Execute(groupID int64) ([]entities.OutingWithDetails, error) {
	return uc.repo.GetByGroupID(groupID)
}
