package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type GetGroup struct {
	repo repository.GroupRepository
}

func NewGetGroup(repo repository.GroupRepository) *GetGroup {
	return &GetGroup{repo: repo}
}

func (uc *GetGroup) Execute(id int64, userID int64) (*entities.Group, error) {
	group, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("grupo no encontrado")
	}

	// Verificar que el usuario sea miembro del grupo
	member, err := uc.repo.GetMemberByGroupAndUser(id, userID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar membresía: %v", err)
	}
	if member == nil || member.Status != entities.MemberStatusAccepted {
		return nil, fmt.Errorf("no tienes acceso a este grupo")
	}

	return group, nil
}
