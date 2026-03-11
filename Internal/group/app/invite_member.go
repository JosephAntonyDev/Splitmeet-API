package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type InviteMember struct {
	groupRepo repository.GroupRepository
	userRepo  userRepository.UserRepository
	notifSvc  *core.NotificationService
}

func NewInviteMember(groupRepo repository.GroupRepository, userRepo userRepository.UserRepository, notifSvc *core.NotificationService) *InviteMember {
	return &InviteMember{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		notifSvc:  notifSvc,
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
		Role:      entities.MemberRoleMember,
		Status:    entities.MemberStatusPending,
		InvitedBy: &input.InviterID,
		InvitedAt: time.Now(),
	}

	err = uc.groupRepo.AddMember(member)
	if err != nil {
		return nil, fmt.Errorf("error al crear invitación: %v", err)
	}

	// Send notification to the invited user
	if uc.notifSvc != nil {
		inviter, _ := uc.userRepo.GetByID(input.InviterID)
		inviterName := "Alguien"
		if inviter != nil {
			inviterName = inviter.Name
		}

		uc.notifSvc.Send(core.NotificationPayload{
			UserID:      user.ID,
			Type:        "group_invitation",
			Title:       "Invitación a grupo",
			Message:     fmt.Sprintf("%s te ha invitado a unirte al grupo %s", inviterName, group.Name),
			ReferenceID: &group.ID,
			InviterName: inviterName,
			GroupName:   group.Name,
		})
	}

	return member, nil
}
