package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type UpdateItemRequest struct {
	CustomName         *string  `json:"custom_name"`
	CustomPresentation *string  `json:"custom_presentation"`
	Quantity           *int     `json:"quantity"`
	UnitPrice          *float64 `json:"unit_price"`
	IsShared           *bool    `json:"is_shared"`
}

type UpdateItemUseCase struct {
	repo repository.OutingRepository
}

func NewUpdateItemUseCase(repo repository.OutingRepository) *UpdateItemUseCase {
	return &UpdateItemUseCase{repo: repo}
}

func (uc *UpdateItemUseCase) Execute(itemID int64, userID int64, req UpdateItemRequest) (*entities.OutingItem, error) {
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

	// Only the outing creator can update items
	if outing.CreatorID != userID {
		return nil, errors.New("only the outing creator can update items")
	}

	if req.CustomName != nil {
		item.CustomName = *req.CustomName
	}
	if req.CustomPresentation != nil {
		item.CustomPresentation = *req.CustomPresentation
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.UnitPrice != nil {
		item.UnitPrice = *req.UnitPrice
	}
	if req.IsShared != nil {
		item.IsShared = *req.IsShared
	}

	item.Subtotal = float64(item.Quantity) * item.UnitPrice

	err = uc.repo.UpdateItem(item)
	if err != nil {
		return nil, err
	}

	// Update outing total
	items, _ := uc.repo.GetItemsByOutingID(item.OutingID)
	var total float64
	for _, i := range items {
		total += i.Subtotal
	}
	uc.repo.UpdateTotalAmount(item.OutingID, total)

	return item, nil
}
