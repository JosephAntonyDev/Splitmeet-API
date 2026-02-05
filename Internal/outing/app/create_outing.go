package app

import (
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type CreateOutingRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CategoryID  *int64 `json:"category_id"`
	GroupID     *int64 `json:"group_id"`
	OutingDate  string `json:"outing_date"`
	SplitType   string `json:"split_type" binding:"required"`
}

type CreateOutingUseCase struct {
	repo repository.OutingRepository
}

func NewCreateOutingUseCase(repo repository.OutingRepository) *CreateOutingUseCase {
	return &CreateOutingUseCase{repo: repo}
}

func (uc *CreateOutingUseCase) Execute(req CreateOutingRequest, creatorID int64) (*entities.Outing, error) {
	var outingDate time.Time
	if req.OutingDate != "" {
		parsed, err := time.Parse("2006-01-02", req.OutingDate)
		if err == nil {
			outingDate = parsed
		}
	}

	outing := &entities.Outing{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		GroupID:     req.GroupID,
		CreatorID:   creatorID,
		OutingDate:  outingDate,
		SplitType:   entities.SplitType(req.SplitType),
		Status:      entities.OutingStatusActive,
		IsEditable:  true,
	}

	err := uc.repo.Save(outing)
	if err != nil {
		return nil, err
	}

	// Add creator as participant automatically
	participant := &entities.OutingParticipant{
		OutingID: outing.ID,
		UserID:   creatorID,
		Status:   entities.ParticipantStatusConfirmed,
	}
	uc.repo.AddParticipant(participant)

	return outing, nil
}
