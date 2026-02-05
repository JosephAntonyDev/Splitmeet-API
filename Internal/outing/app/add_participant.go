package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type AddParticipantRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type AddParticipantUseCase struct {
	repo repository.OutingRepository
}

func NewAddParticipantUseCase(repo repository.OutingRepository) *AddParticipantUseCase {
	return &AddParticipantUseCase{repo: repo}
}

func (uc *AddParticipantUseCase) Execute(outingID int64, inviterID int64, req AddParticipantRequest) (*entities.OutingParticipant, error) {
	// Verify outing exists
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return nil, err
	}

	// Verify inviter is the creator or an existing participant
	if outing.CreatorID != inviterID {
		_, err := uc.repo.GetParticipantByOutingAndUser(outingID, inviterID)
		if err != nil {
			return nil, errors.New("only participants can add other participants")
		}
	}

	// Check if user is already a participant
	existing, _ := uc.repo.GetParticipantByOutingAndUser(outingID, req.UserID)
	if existing != nil {
		return nil, errors.New("user is already a participant")
	}

	participant := &entities.OutingParticipant{
		OutingID: outingID,
		UserID:   req.UserID,
		Status:   entities.ParticipantStatusPending,
	}

	err = uc.repo.AddParticipant(participant)
	if err != nil {
		return nil, err
	}

	return participant, nil
}
