package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type GetItemsUseCase struct {
	repo repository.OutingRepository
}

func NewGetItemsUseCase(repo repository.OutingRepository) *GetItemsUseCase {
	return &GetItemsUseCase{repo: repo}
}

func (uc *GetItemsUseCase) Execute(outingID int64) ([]entities.OutingItemWithProduct, error) {
	return uc.repo.GetItemsByOutingID(outingID)
}
