package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type GetPendingInvitations struct {
	repo repository.GroupRepository
}

func NewGetPendingInvitations(repo repository.GroupRepository) *GetPendingInvitations {
	return &GetPendingInvitations{repo: repo}
}

func (uc *GetPendingInvitations) Execute(userID int64) ([]entities.GroupMember, error) {
	invitations, err := uc.repo.GetPendingInvitations(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener invitaciones: %v", err)
	}
	return invitations, nil
}
