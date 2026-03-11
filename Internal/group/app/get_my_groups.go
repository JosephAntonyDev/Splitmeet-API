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

type GetMyGroupsInput struct {
	UserID int64
	Limit  int
	Offset int
	Search string
}

func (uc *GetMyGroups) Execute(input GetMyGroupsInput) ([]entities.GroupWithDetails, int, error) {
	groups, total, err := uc.repo.GetByUserID(input.UserID, input.Limit, input.Offset, input.Search)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener grupos: %v", err)
	}
	return groups, total, nil
}
