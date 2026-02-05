package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type InviteMember struct {
	groupRepo repository.GroupRepository
	userRepo  userRepository.UserRepository
}

func NewInviteMember(groupRepo repository.GroupRepository, userRepo userRepository.UserRepository) *InviteMember {
	return &InviteMember{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

type InviteMemberInput struct {
	GroupID   int64
	Username  string
	InviterID int64
}

func (uc *InviteMember) Execute(input InviteMemberInput) (*entities.GroupMember, error) {
	// Verificar que el grupo exista
	group, err := uc.groupRepo.GetByID(input.GroupID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("grupo no encontrado")
	}

	// Verificar que quien invita sea miembro aceptado del grupo
	inviterMember, err := uc.groupRepo.GetMemberByGroupAndUser(input.GroupID, input.InviterID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar membresía: %v", err)
	}
	if inviterMember == nil || inviterMember.Status != entities.MemberStatusAccepted {
		return nil, fmt.Errorf("no tienes permisos para invitar a este grupo")
	}

	// Buscar al usuario por username
	user, err := uc.userRepo.GetByUsername(input.Username)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Verificar que no esté ya en el grupo
	existingMember, err := uc.groupRepo.GetMemberByGroupAndUser(input.GroupID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar membresía existente: %v", err)
	}
	if existingMember != nil {
		if existingMember.Status == entities.MemberStatusAccepted {
			return nil, fmt.Errorf("el usuario ya es miembro del grupo")
		}
		if existingMember.Status == entities.MemberStatusPending {
			return nil, fmt.Errorf("el usuario ya tiene una invitación pendiente")
		}
	}

	// Crear la invitación
	member := &entities.GroupMember{
		GroupID:   input.GroupID,
		UserID:    user.ID,
		Status:    entities.MemberStatusPending,
		InvitedBy: &input.InviterID,
		InvitedAt: time.Now(),
	}

	err = uc.groupRepo.AddMember(member)
	if err != nil {
		return nil, fmt.Errorf("error al crear invitación: %v", err)
	}

	return member, nil
}
