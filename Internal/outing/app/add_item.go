package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type AddItemRequest struct {
	ProductID          *int64  `json:"product_id"`
	CustomName         string  `json:"custom_name"`
	CustomPresentation string  `json:"custom_presentation"`
	Quantity           int     `json:"quantity" binding:"required,min=1"`
	UnitPrice          float64 `json:"unit_price" binding:"required"`
	IsShared           bool    `json:"is_shared"`
}

type AddItemUseCase struct {
	repo repository.OutingRepository
}

func NewAddItemUseCase(repo repository.OutingRepository) *AddItemUseCase {
	return &AddItemUseCase{repo: repo}
}

func (uc *AddItemUseCase) Execute(outingID int64, userID int64, req AddItemRequest) (*entities.OutingItem, error) {
	// Verify outing exists
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return nil, err
	}

	if !outing.IsEditable {
		return nil, errors.New("outing is no longer editable")
	}

	// Verify user is a participant
	_, err = uc.repo.GetParticipantByOutingAndUser(outingID, userID)
	if err != nil {
		return nil, errors.New("only participants can add items")
	}

	subtotal := float64(req.Quantity) * req.UnitPrice

	item := &entities.OutingItem{
		OutingID:           outingID,
		ProductID:          req.ProductID,
		CustomName:         req.CustomName,
		CustomPresentation: req.CustomPresentation,
		Quantity:           req.Quantity,
		UnitPrice:          req.UnitPrice,
		Subtotal:           subtotal,
		IsShared:           req.IsShared,
	}

	err = uc.repo.AddItem(item)
	if err != nil {
		return nil, err
	}

	// Update outing total
	items, _ := uc.repo.GetItemsByOutingID(outingID)
	var total float64
	for _, i := range items {
		total += i.Subtotal
	}
	uc.repo.UpdateTotalAmount(outingID, total)

	return item, nil
}
