package app

import (
	"errors"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type JoinGroupUseCase struct {
	repo repository.GroupRepository
}

func NewJoinGroupUseCase(repo repository.GroupRepository) *JoinGroupUseCase {
	return &JoinGroupUseCase{repo: repo}
}

func (uc *JoinGroupUseCase) Execute(groupID int64, userID int64) (*entities.GroupMember, error) {
	// Verificar si el grupo existe
	group, err := uc.repo.GetByID(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	// Verificar si ya es miembro o está pendiente
	existing, _ := uc.repo.GetMemberByGroupAndUser(groupID, userID)
	if existing != nil {
		if existing.Status == entities.MemberStatusAccepted {
			return nil, errors.New("user is already a member of this group")
		} else if existing.Status == entities.MemberStatusPending {
			// Podría actualizar aquí mismo si estaba pendiente, o rechazarlo
			return nil, errors.New("user already has a pending invitation, please respond to it")
		}
	}

	// Agregar directamente al grupo
	now := time.Now()
	member := &entities.GroupMember{
		GroupID:     groupID,
		UserID:      userID,
		Role:        entities.MemberRoleMember,
		Status:      entities.MemberStatusAccepted,
		InvitedBy:   &userID, // self-invited
		InvitedAt:   now,
		RespondedAt: &now,
	}

	err = uc.repo.AddMember(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
