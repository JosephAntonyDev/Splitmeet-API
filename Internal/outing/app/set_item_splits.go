package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type SplitDefinition struct {
	ParticipantID int64    `json:"participant_id" binding:"required"`
	Amount        float64  `json:"amount" binding:"required"`
	Percentage    *float64 `json:"percentage"`
}

type SetItemSplitsRequest struct {
	Splits []SplitDefinition `json:"splits" binding:"required"`
}

type SetItemSplitsUseCase struct {
	repo repository.OutingRepository
}

func NewSetItemSplitsUseCase(repo repository.OutingRepository) *SetItemSplitsUseCase {
	return &SetItemSplitsUseCase{repo: repo}
}

func (uc *SetItemSplitsUseCase) Execute(itemID int64, userID int64, req SetItemSplitsRequest) ([]entities.ItemSplitWithUser, error) {
	item, err := uc.repo.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}

	outing, err := uc.repo.GetByID(item.OutingID)
	if err != nil {
		return nil, err
	}

	if !outing.IsEditable {
		return nil, errors.New("outing is no longer editable")
	}

	// Only outing creator can set splits
	if outing.CreatorID != userID {
		return nil, errors.New("only the outing creator can set splits")
	}

	// Validate total amount equals item subtotal
	var totalAmount float64
	for _, s := range req.Splits {
		totalAmount += s.Amount
	}

	if totalAmount != item.Subtotal {
		return nil, errors.New("split amounts must equal item subtotal")
	}

	// Remove existing splits
	uc.repo.DeleteSplitsByItemID(itemID)

	// Create new splits
	for _, s := range req.Splits {
		split := &entities.ItemSplit{
			OutingItemID:  itemID,
			ParticipantID: s.ParticipantID,
			SplitAmount:   s.Amount,
			Percentage:    s.Percentage,
		}
		uc.repo.AddItemSplit(split)
	}

	return uc.repo.GetSplitsByItemID(itemID)
}
