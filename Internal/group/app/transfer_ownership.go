package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type TransferOwnership struct {
	repo repository.GroupRepository
}

func NewTransferOwnership(repo repository.GroupRepository) *TransferOwnership {
	return &TransferOwnership{repo: repo}
}

type TransferOwnershipInput struct {
	GroupID     int64
	NewOwnerID  int64
	RequestedBy int64
}

func (uc *TransferOwnership) Execute(input TransferOwnershipInput) error {
	group, err := uc.repo.GetByID(input.GroupID)
	if err != nil {
		return fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return fmt.Errorf("grupo no encontrado")
	}

	if group.OwnerID != input.RequestedBy {
		return fmt.Errorf("solo el propietario puede transferir la propiedad del grupo")
	}

	// Verify new owner is an accepted member
	newOwnerMember, err := uc.repo.GetMemberByGroupAndUser(input.GroupID, input.NewOwnerID)
	if err != nil {
		return fmt.Errorf("error al verificar miembro: %v", err)
	}
	if newOwnerMember == nil || newOwnerMember.Status != entities.MemberStatusAccepted {
		return fmt.Errorf("el nuevo propietario debe ser un miembro activo del grupo")
	}

	return uc.repo.TransferOwnership(input.GroupID, input.NewOwnerID)
}
