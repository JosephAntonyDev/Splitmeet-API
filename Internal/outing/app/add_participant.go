package app

import (
	"errors"
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type AddParticipantRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type AddParticipantUseCase struct {
	repo     repository.OutingRepository
	userRepo userRepository.UserRepository
	notifSvc *core.NotificationService
}

func NewAddParticipantUseCase(repo repository.OutingRepository, userRepo userRepository.UserRepository, notifSvc *core.NotificationService) *AddParticipantUseCase {
	return &AddParticipantUseCase{repo: repo, userRepo: userRepo, notifSvc: notifSvc}
}

func (uc *AddParticipantUseCase) Execute(outingID int64, inviterID int64, req AddParticipantRequest) (*entities.OutingParticipant, error) {
	// Verify outing exists
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return nil, err
	}

	// Verify inviter is the creator or an existing participant
	if outing.CreatorID != inviterID {
		_, err := uc.repo.GetParticipantByOutingAndUser(outingID, inviterID)
		if err != nil {
			return nil, errors.New("only participants can add other participants")
		}
	}

	// Check if user is already a participant
	existing, _ := uc.repo.GetParticipantByOutingAndUser(outingID, req.UserID)
	if existing != nil {
		return nil, errors.New("user is already a participant")
	}

	participant := &entities.OutingParticipant{
		OutingID:  outingID,
		UserID:    req.UserID,
		InvitedBy: &inviterID,
		Status:    entities.ParticipantStatusPending,
	}

	err = uc.repo.AddParticipant(participant)
	if err != nil {
		return nil, err
	}

	// Send notification to invited user
	if uc.notifSvc != nil {
		inviter, _ := uc.userRepo.GetByID(inviterID)
		inviterName := ""
		if inviter != nil {
			inviterName = inviter.Name
		}
		uc.notifSvc.Send(core.NotificationPayload{
			UserID:      req.UserID,
			Type:        "outing_invitation",
			Title:       "Invitación a salida",
			Message:     fmt.Sprintf("%s te invitó a la salida %s", inviterName, outing.Name),
			ReferenceID: &outing.ID,
			InviterName: inviterName,
			OutingName:  outing.Name,
		})
	}

	return participant, nil
}
