package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type CreateGroup struct {
	repo repository.GroupRepository
}

func NewCreateGroup(repo repository.GroupRepository) *CreateGroup {
	return &CreateGroup{repo: repo}
}

type CreateGroupInput struct {
	Name        string
	Description string
	OwnerID     int64
}

func (uc *CreateGroup) Execute(input CreateGroupInput) (*entities.Group, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("el nombre del grupo es requerido")
	}

	now := time.Now()
	group := &entities.Group{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     input.OwnerID,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := uc.repo.Save(group)
	if err != nil {
		return nil, fmt.Errorf("error al crear grupo: %v", err)
	}

	// Agregar al owner como miembro aceptado automáticamente
	member := &entities.GroupMember{
		GroupID:   group.ID,
		UserID:    input.OwnerID,
		Role:      entities.MemberRoleOwner,
		Status:    entities.MemberStatusAccepted,
		InvitedBy: &input.OwnerID,
		InvitedAt: now,
	}
	respondedAt := now
	member.RespondedAt = &respondedAt

	err = uc.repo.AddMember(member)
	if err != nil {
		return nil, fmt.Errorf("error al agregar owner como miembro: %v", err)
	}

	return group, nil
}
