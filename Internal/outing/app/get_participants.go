package app

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type GetParticipantsUseCase struct {
	repo repository.OutingRepository
}

func NewGetParticipantsUseCase(repo repository.OutingRepository) *GetParticipantsUseCase {
	return &GetParticipantsUseCase{repo: repo}
}

func (uc *GetParticipantsUseCase) Execute(outingID int64) ([]entities.OutingParticipantWithUser, error) {
	return uc.repo.GetParticipantsByOutingID(outingID)
}
