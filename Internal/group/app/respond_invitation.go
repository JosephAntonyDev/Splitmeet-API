package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type RespondInvitation struct {
	repo     repository.GroupRepository
	userRepo userRepository.UserRepository
	notifSvc *core.NotificationService
}

func NewRespondInvitation(repo repository.GroupRepository, userRepo userRepository.UserRepository, notifSvc *core.NotificationService) *RespondInvitation {
	return &RespondInvitation{repo: repo, userRepo: userRepo, notifSvc: notifSvc}
}

type RespondInvitationInput struct {
	GroupID int64
	UserID  int64
	Accept  bool
}

func (uc *RespondInvitation) Execute(input RespondInvitationInput) error {
	// Verificar que exista la invitación
	member, err := uc.repo.GetMemberByGroupAndUser(input.GroupID, input.UserID)
	if err != nil {
		return fmt.Errorf("error al buscar invitación: %v", err)
	}
	if member == nil {
		return fmt.Errorf("no tienes una invitación a este grupo")
	}
	if member.Status != entities.MemberStatusPending {
		return fmt.Errorf("la invitación ya fue respondida")
	}

	var newStatus entities.MemberStatus
	if input.Accept {
		newStatus = entities.MemberStatusAccepted
	} else {
		newStatus = entities.MemberStatusRejected
	}

	err = uc.repo.UpdateMemberStatus(input.GroupID, input.UserID, newStatus)
	if err != nil {
		return fmt.Errorf("error al actualizar invitación: %v", err)
	}

	// Send notification to the group owner about the response
	if uc.notifSvc != nil {
		group, _ := uc.repo.GetByID(input.GroupID)
		responder, _ := uc.userRepo.GetByID(input.UserID)

		if group != nil && responder != nil {
			notifType := "invitation_accepted"
			action := "aceptó"
			if !input.Accept {
				notifType = "invitation_rejected"
				action = "rechazó"
			}

			uc.notifSvc.Send(core.NotificationPayload{
				UserID:      group.OwnerID,
				Type:        notifType,
				Title:       "Respuesta a invitación",
				Message:     fmt.Sprintf("%s %s la invitación al grupo %s", responder.Name, action, group.Name),
				ReferenceID: &group.ID,
				InviterName: responder.Name,
				GroupName:   group.Name,
			})
		}
	}

	return nil
}
