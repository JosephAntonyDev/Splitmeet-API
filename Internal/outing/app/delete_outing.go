package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type DeleteOutingUseCase struct {
	repo repository.OutingRepository
}

func NewDeleteOutingUseCase(repo repository.OutingRepository) *DeleteOutingUseCase {
	return &DeleteOutingUseCase{repo: repo}
}

func (uc *DeleteOutingUseCase) Execute(outingID int64, userID int64) error {
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return err
	}

	if outing.CreatorID != userID {
		return errors.New("only the creator can delete the outing")
	}

	return uc.repo.Delete(outingID)
}
