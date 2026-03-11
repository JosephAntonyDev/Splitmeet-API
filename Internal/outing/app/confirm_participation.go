package app

import (
	"errors"
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type ConfirmParticipationUseCase struct {
	repo     repository.OutingRepository
	userRepo userRepository.UserRepository
	notifSvc *core.NotificationService
}

func NewConfirmParticipationUseCase(repo repository.OutingRepository, userRepo userRepository.UserRepository, notifSvc *core.NotificationService) *ConfirmParticipationUseCase {
	return &ConfirmParticipationUseCase{repo: repo, userRepo: userRepo, notifSvc: notifSvc}
}

func (uc *ConfirmParticipationUseCase) Execute(outingID int64, userID int64, accept bool) error {
	participant, err := uc.repo.GetParticipantByOutingAndUser(outingID, userID)
	if err != nil {
		return err
	}

	if participant.Status != entities.ParticipantStatusPending {
		return errors.New("participation already confirmed or declined")
	}

	var status entities.ParticipantStatus
	if accept {
		status = entities.ParticipantStatusConfirmed
	} else {
		status = entities.ParticipantStatusDeclined
	}

	err = uc.repo.UpdateParticipantStatus(outingID, userID, status)
	if err != nil {
		return err
	}

	// Notify outing creator about the response
	if uc.notifSvc != nil {
		outing, _ := uc.repo.GetByID(outingID)
		responder, _ := uc.userRepo.GetByID(userID)

		if outing != nil && responder != nil {
			notifType := "invitation_accepted"
			action := "aceptó"
			if !accept {
				notifType = "invitation_rejected"
				action = "rechazó"
			}

			uc.notifSvc.Send(core.NotificationPayload{
				UserID:      outing.CreatorID,
				Type:        notifType,
				Title:       "Respuesta a invitación de salida",
				Message:     fmt.Sprintf("%s %s la invitación a la salida %s", responder.Name, action, outing.Name),
				ReferenceID: &outing.ID,
				InviterName: responder.Name,
				OutingName:  outing.Name,
			})
		}
	}

	return nil
}
