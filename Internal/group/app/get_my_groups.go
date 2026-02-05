package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type GetMyGroups struct {
	repo repository.GroupRepository
}

func NewGetMyGroups(repo repository.GroupRepository) *GetMyGroups {
	return &GetMyGroups{repo: repo}
}

func (uc *GetMyGroups) Execute(userID int64) ([]entities.Group, error) {
	groups, err := uc.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener grupos: %v", err)
	}
	return groups, nil
}
