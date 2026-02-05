package app

import (
	"fmt"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type DeleteUser struct {
	repo repository.UserRepository
}

func NewDeleteUser(repo repository.UserRepository) *DeleteUser {
	return &DeleteUser{repo: repo}
}

func (uc *DeleteUser) Execute(id int64) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	return nil
}