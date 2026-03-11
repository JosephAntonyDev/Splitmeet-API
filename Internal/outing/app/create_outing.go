package app

import (
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	groupRepository "github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type CreateOutingRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CategoryID  *int64 `json:"category_id"`
	GroupID     *int64 `json:"group_id"`
	OutingDate  string `json:"outing_date"`
	SplitType   string `json:"split_type" binding:"required"`
}

type CreateOutingUseCase struct {
	repo      repository.OutingRepository
	groupRepo groupRepository.GroupRepository
	userRepo  userRepository.UserRepository
	notifSvc  *core.NotificationService
}

func NewCreateOutingUseCase(repo repository.OutingRepository, groupRepo groupRepository.GroupRepository, userRepo userRepository.UserRepository, notifSvc *core.NotificationService) *CreateOutingUseCase {
	return &CreateOutingUseCase{repo: repo, groupRepo: groupRepo, userRepo: userRepo, notifSvc: notifSvc}
}

func (uc *CreateOutingUseCase) Execute(req CreateOutingRequest, creatorID int64) (*entities.Outing, error) {
	var outingDate time.Time
	if req.OutingDate != "" {
		parsed, err := time.Parse("2006-01-02", req.OutingDate)
		if err == nil {
			outingDate = parsed
		}
	}

	outing := &entities.Outing{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		GroupID:     req.GroupID,
		CreatorID:   creatorID,
		OutingDate:  outingDate,
		SplitType:   entities.SplitType(req.SplitType),
		Status:      entities.OutingStatusActive,
		IsEditable:  true,
	}

	err := uc.repo.Save(outing)
	if err != nil {
		return nil, err
	}

	// Add creator as participant automatically
	participant := &entities.OutingParticipant{
		OutingID: outing.ID,
		UserID:   creatorID,
		Status:   entities.ParticipantStatusConfirmed,
	}
	uc.repo.AddParticipant(participant)

	// Auto-invite group members if outing is linked to a group
	if req.GroupID != nil && uc.groupRepo != nil {
		memberIDs, err := uc.groupRepo.GetAcceptedMemberIDs(*req.GroupID)
		if err == nil {
			creator, _ := uc.userRepo.GetByID(creatorID)
			creatorName := ""
			if creator != nil {
				creatorName = creator.Name
			}
			group, _ := uc.groupRepo.GetByID(*req.GroupID)
			groupName := ""
			if group != nil {
				groupName = group.Name
			}

			for _, memberID := range memberIDs {
				if memberID == creatorID {
					continue
				}
				p := &entities.OutingParticipant{
					OutingID:  outing.ID,
					UserID:    memberID,
					InvitedBy: &creatorID,
					Status:    entities.ParticipantStatusPending,
				}
				uc.repo.AddParticipant(p)

				if uc.notifSvc != nil {
					uc.notifSvc.Send(core.NotificationPayload{
						UserID:      memberID,
						Type:        "outing_invitation",
						Title:       "Invitación a salida",
						Message:     fmt.Sprintf("%s te invitó a la salida %s", creatorName, outing.Name),
						ReferenceID: &outing.ID,
						InviterName: creatorName,
						GroupName:   groupName,
						OutingName:  outing.Name,
					})
				}
			}
		}
	}

	return outing, nil
}
