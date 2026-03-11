package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type RemoveMember struct {
	repo repository.GroupRepository
}

func NewRemoveMember(repo repository.GroupRepository) *RemoveMember {
	return &RemoveMember{repo: repo}
}

type RemoveMemberInput struct {
	GroupID        int64
	MemberToRemove int64
	RequestedBy    int64
}

func (uc *RemoveMember) Execute(input RemoveMemberInput) error {
	// Obtener el grupo
	group, err := uc.repo.GetByID(input.GroupID)
	if err != nil {
		return fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return fmt.Errorf("grupo no encontrado")
	}

	// Verificar que quien solicita sea miembro aceptado
	requester, err := uc.repo.GetMemberByGroupAndUser(input.GroupID, input.RequestedBy)
	if err != nil {
		return fmt.Errorf("error al verificar membresía: %v", err)
	}
	if requester == nil || requester.Status != entities.MemberStatusAccepted {
		return fmt.Errorf("no tienes acceso a este grupo")
	}

	// El owner y admins pueden remover a otros, o un usuario puede salirse a sí mismo
	if input.MemberToRemove != input.RequestedBy {
		if requester.Role != entities.MemberRoleOwner && requester.Role != entities.MemberRoleAdmin {
			return fmt.Errorf("solo el propietario o administradores pueden remover miembros")
		}
	}

	// No se puede remover al owner
	if input.MemberToRemove == group.OwnerID && input.RequestedBy != group.OwnerID {
		return fmt.Errorf("no se puede remover al propietario del grupo")
	}

	// Si el owner quiere salirse, debe transferir o eliminar el grupo
	if input.MemberToRemove == group.OwnerID {
		return fmt.Errorf("el propietario no puede abandonar el grupo, debe transferir la propiedad o eliminarlo")
	}

	// Un admin no puede remover a otro admin (solo el owner puede)
	if input.MemberToRemove != input.RequestedBy {
		target, err := uc.repo.GetMemberByGroupAndUser(input.GroupID, input.MemberToRemove)
		if err != nil {
			return fmt.Errorf("error al verificar miembro: %v", err)
		}
		if target != nil && target.Role == entities.MemberRoleAdmin && requester.Role != entities.MemberRoleOwner {
			return fmt.Errorf("solo el propietario puede remover administradores")
		}
	}

	err = uc.repo.RemoveMember(input.GroupID, input.MemberToRemove)
	if err != nil {
		return fmt.Errorf("error al remover miembro: %v", err)
	}

	return nil
}
