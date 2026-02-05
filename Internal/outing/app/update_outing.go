package app

import (
	"errors"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type UpdateOutingRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	CategoryID  *int64  `json:"category_id"`
	OutingDate  *string `json:"outing_date"`
	SplitType   *string `json:"split_type"`
	Status      *string `json:"status"`
}

type UpdateOutingUseCase struct {
	repo repository.OutingRepository
}

func NewUpdateOutingUseCase(repo repository.OutingRepository) *UpdateOutingUseCase {
	return &UpdateOutingUseCase{repo: repo}
}

func (uc *UpdateOutingUseCase) Execute(outingID int64, userID int64, req UpdateOutingRequest) (*entities.Outing, error) {
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return nil, err
	}

	if outing.CreatorID != userID {
		return nil, errors.New("only the creator can update the outing")
	}

	if !outing.IsEditable {
		return nil, errors.New("outing is no longer editable")
	}

	if req.Name != nil {
		outing.Name = *req.Name
	}
	if req.Description != nil {
		outing.Description = *req.Description
	}
	if req.CategoryID != nil {
		outing.CategoryID = req.CategoryID
	}
	if req.OutingDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.OutingDate)
		if err == nil {
			outing.OutingDate = parsed
		}
	}
	if req.SplitType != nil {
		outing.SplitType = entities.SplitType(*req.SplitType)
	}
	if req.Status != nil {
		outing.Status = entities.OutingStatus(*req.Status)
	}

	err = uc.repo.Update(outing)
	if err != nil {
		return nil, err
	}

	return outing, nil
}
