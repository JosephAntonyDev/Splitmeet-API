package app

import (
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/repository"
)

type DeleteGroup struct {
	repo repository.GroupRepository
}

func NewDeleteGroup(repo repository.GroupRepository) *DeleteGroup {
	return &DeleteGroup{repo: repo}
}

func (uc *DeleteGroup) Execute(groupID int64, userID int64) error {
	group, err := uc.repo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("error al buscar grupo: %v", err)
	}
	if group == nil {
		return fmt.Errorf("grupo no encontrado")
	}

	// Solo el owner puede eliminar el grupo
	if group.OwnerID != userID {
		return fmt.Errorf("solo el creador del grupo puede eliminarlo")
	}

	err = uc.repo.Delete(groupID)
	if err != nil {
		return fmt.Errorf("error al eliminar grupo: %v", err)
	}

	return nil
}
