package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type GetItemSplitsUseCase struct {
	repo repository.OutingRepository
}

func NewGetItemSplitsUseCase(repo repository.OutingRepository) *GetItemSplitsUseCase {
	return &GetItemSplitsUseCase{repo: repo}
}

func (uc *GetItemSplitsUseCase) Execute(itemID int64) ([]entities.ItemSplitWithUser, error) {
	return uc.repo.GetSplitsByItemID(itemID)
}
