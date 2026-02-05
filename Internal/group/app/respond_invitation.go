package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type RespondInvitation struct {
	repo repository.GroupRepository
}

func NewRespondInvitation(repo repository.GroupRepository) *RespondInvitation {
	return &RespondInvitation{repo: repo}
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

	// Actualizar responded_at (esto se puede hacer en el repo también)
	_ = time.Now() // timestamp de respuesta

	return nil
}
