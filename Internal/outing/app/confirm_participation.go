package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type ConfirmParticipationUseCase struct {
	repo repository.OutingRepository
}

func NewConfirmParticipationUseCase(repo repository.OutingRepository) *ConfirmParticipationUseCase {
	return &ConfirmParticipationUseCase{repo: repo}
}

func (uc *ConfirmParticipationUseCase) Execute(outingID int64, userID int64, accept bool) error {
	participant, err := uc.repo.GetParticipantByOutingAndUser(outingID, userID)
	if err != nil {
		return err
	}

	if participant.Status != entities.ParticipantStatusPending {
		return errors.New("participation already confirmed or declined")
	}

	var status entities.ParticipantStatus
	if accept {
		status = entities.ParticipantStatusConfirmed
	} else {
		status = entities.ParticipantStatusDeclined
	}

	return uc.repo.UpdateParticipantStatus(outingID, userID, status)
}
