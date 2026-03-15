package app

import (
	"errors"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type JoinOutingUseCase struct {
	repo repository.OutingRepository
}

func NewJoinOutingUseCase(repo repository.OutingRepository) *JoinOutingUseCase {
	return &JoinOutingUseCase{repo: repo}
}

func (uc *JoinOutingUseCase) Execute(outingID int64, userID int64) (*entities.OutingParticipant, error) {
	// 1. Verify outing exists
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return nil, err
	}
	if outing == nil {
		return nil, errors.New("outing not found")
	}

	// 2. Check if user is already a participant
	existing, _ := uc.repo.GetParticipantByOutingAndUser(outingID, userID)
	if existing != nil {
		return nil, errors.New("user is already a participant")
	}

	// 3. Add participant directly as "confirmed" (since they join voluntarily via QR)
	now := time.Now()
	participant := &entities.OutingParticipant{
		OutingID:  outingID,
		UserID:    userID,
		InvitedBy: &userID, // self-invited
		Status:    entities.ParticipantStatusConfirmed,
		JoinedAt:  now,
	}

	err = uc.repo.AddParticipant(participant)
	if err != nil {
		return nil, err
	}

	return participant, nil
}
