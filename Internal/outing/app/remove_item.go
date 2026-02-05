package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type RemoveItemUseCase struct {
	repo repository.OutingRepository
}

func NewRemoveItemUseCase(repo repository.OutingRepository) *RemoveItemUseCase {
	return &RemoveItemUseCase{repo: repo}
}

func (uc *RemoveItemUseCase) Execute(itemID int64, userID int64) error {
	item, err := uc.repo.GetItemByID(itemID)
	if err != nil {
		return err
	}

	outing, err := uc.repo.GetByID(item.OutingID)
	if err != nil {
		return err
	}

	if !outing.IsEditable {
		return errors.New("outing is no longer editable")
	}

	// Only the outing creator can remove items
	if outing.CreatorID != userID {
		return errors.New("only the outing creator can remove items")
	}

	// Delete splits first
	uc.repo.DeleteSplitsByItemID(itemID)

	err = uc.repo.DeleteItem(itemID)
	if err != nil {
		return err
	}

	// Update outing total
	items, _ := uc.repo.GetItemsByOutingID(item.OutingID)
	var total float64
	for _, i := range items {
		total += i.Subtotal
	}
	uc.repo.UpdateTotalAmount(item.OutingID, total)

	return nil
}
