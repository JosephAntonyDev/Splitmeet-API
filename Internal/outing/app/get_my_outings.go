package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type GetMyOutingsUseCase struct {
	repo repository.OutingRepository
}

func NewGetMyOutingsUseCase(repo repository.OutingRepository) *GetMyOutingsUseCase {
	return &GetMyOutingsUseCase{repo: repo}
}

func (uc *GetMyOutingsUseCase) Execute(userID int64) ([]entities.OutingWithDetails, error) {
	return uc.repo.GetByUserID(userID)
}
