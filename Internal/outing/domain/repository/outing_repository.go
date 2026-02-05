package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"

type OutingRepository interface {
	// Outing operations
	Save(outing *entities.Outing) error
	GetByID(id int64) (*entities.Outing, error)
	GetByIDWithDetails(id int64) (*entities.OutingWithDetails, error)
	GetByUserID(userID int64) ([]entities.OutingWithDetails, error)
	GetByGroupID(groupID int64) ([]entities.OutingWithDetails, error)
	Update(outing *entities.Outing) error
	Delete(id int64) error
	UpdateTotalAmount(outingID int64, amount float64) error
	MarkAsCompleted(outingID int64) error

	// Participant operations
	AddParticipant(participant *entities.OutingParticipant) error
	GetParticipantByOutingAndUser(outingID, userID int64) (*entities.OutingParticipant, error)
	GetParticipantByID(participantID int64) (*entities.OutingParticipant, error)
	GetParticipantsByOutingID(outingID int64) ([]entities.OutingParticipantWithUser, error)
	GetConfirmedParticipants(outingID int64) ([]entities.OutingParticipantWithUser, error)
	UpdateParticipantStatus(outingID, userID int64, status entities.ParticipantStatus) error
	UpdateParticipantAmountOwed(participantID int64, amount float64) error
	RemoveParticipant(outingID, userID int64) error

	// Item operations
	AddItem(item *entities.OutingItem) error
	GetItemByID(itemID int64) (*entities.OutingItem, error)
	GetItemsByOutingID(outingID int64) ([]entities.OutingItemWithProduct, error)
	UpdateItem(item *entities.OutingItem) error
	DeleteItem(itemID int64) error

	// Split operations
	AddItemSplit(split *entities.ItemSplit) error
	GetSplitsByItemID(itemID int64) ([]entities.ItemSplitWithUser, error)
	GetSplitsByParticipantID(participantID int64) ([]entities.ItemSplit, error)
	DeleteSplitsByItemID(itemID int64) error
	DeleteSplit(splitID int64) error
}
