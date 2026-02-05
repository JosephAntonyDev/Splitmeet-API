package repository

import "github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"

type GroupRepository interface {
	// Group operations
	Save(group *entities.Group) error
	GetByID(id int64) (*entities.Group, error)
	GetByOwnerID(ownerID int64) ([]entities.Group, error)
	GetByUserID(userID int64) ([]entities.Group, error)
	Update(group *entities.Group) error
	Delete(id int64) error

	// Member operations
	AddMember(member *entities.GroupMember) error
	GetMemberByGroupAndUser(groupID, userID int64) (*entities.GroupMember, error)
	GetMembersByGroupID(groupID int64) ([]entities.GroupMemberWithUser, error)
	UpdateMemberStatus(groupID, userID int64, status entities.MemberStatus) error
	RemoveMember(groupID, userID int64) error
	GetPendingInvitations(userID int64) ([]entities.GroupMember, error)
}
