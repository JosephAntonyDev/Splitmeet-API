package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type SetMemberRole struct {
	repo repository.GroupRepository
}

func NewSetMemberRole(repo repository.GroupRepository) *SetMemberRole {
	return &SetMemberRole{repo: repo}
}

type SetMemberRoleInput struct {
	GroupID      int64
	TargetUserID int64
	Role         string
	RequestedBy  int64
}

func (uc *SetMemberRole) Execute(input SetMemberRoleInput) error {
	group, err := uc.repo.GetByID(input.GroupID)
	if err != nil {
		return fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return fmt.Errorf("grupo no encontrado")
	}

	// Only owner can change roles
	if group.OwnerID != input.RequestedBy {
		return fmt.Errorf("solo el propietario puede cambiar roles")
	}

	// Can't change owner's own role
	if input.TargetUserID == input.RequestedBy {
		return fmt.Errorf("no puedes cambiar tu propio rol, usa transferir propiedad")
	}

	// Validate target is accepted member
	target, err := uc.repo.GetMemberByGroupAndUser(input.GroupID, input.TargetUserID)
	if err != nil {
		return fmt.Errorf("error al verificar miembro: %v", err)
	}
	if target == nil || target.Status != entities.MemberStatusAccepted {
		return fmt.Errorf("el usuario debe ser un miembro activo del grupo")
	}

	// Validate role
	var newRole entities.MemberRole
	switch input.Role {
	case "admin":
		newRole = entities.MemberRoleAdmin
	case "member":
		newRole = entities.MemberRoleMember
	default:
		return fmt.Errorf("rol inválido, debe ser 'admin' o 'member'")
	}

	return uc.repo.UpdateMemberRole(input.GroupID, input.TargetUserID, newRole)
}
