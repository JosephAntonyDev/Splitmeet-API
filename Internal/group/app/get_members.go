package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type GetMembers struct {
	repo repository.GroupRepository
}

func NewGetMembers(repo repository.GroupRepository) *GetMembers {
	return &GetMembers{repo: repo}
}

func (uc *GetMembers) Execute(groupID int64, userID int64) ([]entities.GroupMemberWithUser, error) {
	// Verificar que el usuario sea miembro del grupo
	member, err := uc.repo.GetMemberByGroupAndUser(groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar membresía: %v", err)
	}
	if member == nil || member.Status != entities.MemberStatusAccepted {
		return nil, fmt.Errorf("no tienes acceso a este grupo")
	}

	members, err := uc.repo.GetMembersByGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener miembros: %v", err)
	}

	return members, nil
}
