package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type RemoveParticipantUseCase struct {
	repo repository.OutingRepository
}

func NewRemoveParticipantUseCase(repo repository.OutingRepository) *RemoveParticipantUseCase {
	return &RemoveParticipantUseCase{repo: repo}
}

func (uc *RemoveParticipantUseCase) Execute(outingID int64, userID int64, removerID int64) error {
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return err
	}

	// User can remove themselves, or creator can remove anyone
	if userID != removerID && outing.CreatorID != removerID {
		return errors.New("only the creator can remove other participants")
	}

	return uc.repo.RemoveParticipant(outingID, userID)
}
