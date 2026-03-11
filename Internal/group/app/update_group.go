package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type UpdateGroup struct {
	repo repository.GroupRepository
}

func NewUpdateGroup(repo repository.GroupRepository) *UpdateGroup {
	return &UpdateGroup{repo: repo}
}

type UpdateGroupInput struct {
	GroupID     int64
	Name        string
	Description string
	UserID      int64
}

func (uc *UpdateGroup) Execute(input UpdateGroupInput) (*entities.Group, error) {
	group, err := uc.repo.GetByID(input.GroupID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("grupo no encontrado")
	}

	// Only owner or admin can update
	member, err := uc.repo.GetMemberByGroupAndUser(input.GroupID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar membresía: %v", err)
	}
	if member == nil || (member.Role != entities.MemberRoleOwner && member.Role != entities.MemberRoleAdmin) {
		return nil, fmt.Errorf("solo el propietario o administradores pueden editar el grupo")
	}

	if input.Name != "" {
		group.Name = input.Name
	}
	if input.Description != "" {
		group.Description = input.Description
	}
	group.UpdatedAt = time.Now()

	err = uc.repo.Update(group)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar grupo: %v", err)
	}

	return group, nil
}
