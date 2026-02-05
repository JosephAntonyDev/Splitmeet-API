package app

import (
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/repository"
)

type UserSummary struct {
	ParticipantID int64   `json:"participant_id"`
	UserID        int64   `json:"user_id"`
	Username      string  `json:"username"`
	Name          string  `json:"name"`
	TotalOwed     float64 `json:"total_owed"`
	ItemsCount    int     `json:"items_count"`
}

type CalculateSplitsResponse struct {
	OutingID      int64              `json:"outing_id"`
	TotalAmount   float64            `json:"total_amount"`
	SplitType     entities.SplitType `json:"split_type"`
	UserSummaries []UserSummary      `json:"user_summaries"`
}

type CalculateSplitsUseCase struct {
	repo repository.OutingRepository
}

func NewCalculateSplitsUseCase(repo repository.OutingRepository) *CalculateSplitsUseCase {
	return &CalculateSplitsUseCase{repo: repo}
}

func (uc *CalculateSplitsUseCase) Execute(outingID int64) (*CalculateSplitsResponse, error) {
	outing, err := uc.repo.GetByID(outingID)
	if err != nil {
		return nil, err
	}

	items, err := uc.repo.GetItemsByOutingID(outingID)
	if err != nil {
		return nil, err
	}

	participants, err := uc.repo.GetConfirmedParticipants(outingID)
	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return nil, errors.New("no confirmed participants")
	}

	// Build user summaries map
	summaryMap := make(map[int64]*UserSummary)
	for _, p := range participants {
		summaryMap[p.ID] = &UserSummary{
			ParticipantID: p.ID,
			UserID:        p.UserID,
			Username:      p.Username,
			Name:          p.Name,
			TotalOwed:     0,
			ItemsCount:    0,
		}
	}

	switch outing.SplitType {
	case entities.SplitTypeEqual:
		// Split equally among all confirmed participants
		perPerson := outing.TotalAmount / float64(len(participants))
		for _, summary := range summaryMap {
			summary.TotalOwed = perPerson
			summary.ItemsCount = len(items)
		}

	case entities.SplitTypePerConsumption, entities.SplitTypeCustomFixed:
		// Use existing item splits
		for _, item := range items {
			splits, err := uc.repo.GetSplitsByItemID(item.ID)
			if err != nil {
				continue
			}
			for _, split := range splits {
				if summary, ok := summaryMap[split.ParticipantID]; ok {
					summary.TotalOwed += split.SplitAmount
					summary.ItemsCount++
				}
			}
		}

	case entities.SplitTypeSinglePayer:
		// All goes to the outing creator
		for _, summary := range summaryMap {
			if summary.UserID == outing.CreatorID {
				summary.TotalOwed = outing.TotalAmount
				summary.ItemsCount = len(items)
			}
		}
	}

	// Update participant amounts in database
	for _, summary := range summaryMap {
		uc.repo.UpdateParticipantAmountOwed(summary.ParticipantID, summary.TotalOwed)
	}

	// Convert map to slice
	var summaries []UserSummary
	for _, summary := range summaryMap {
		summaries = append(summaries, *summary)
	}

	return &CalculateSplitsResponse{
		OutingID:      outingID,
		TotalAmount:   outing.TotalAmount,
		SplitType:     outing.SplitType,
		UserSummaries: summaries,
	}, nil
}
