package app

import (
	"fmt"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/ports"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/repository"
)

type UpdateUser struct {
	repo   repository.UserRepository
	bcrypt ports.IBcryptService
}

func NewUpdateUser(repo repository.UserRepository, bcrypt ports.IBcryptService) *UpdateUser {
	return &UpdateUser{
		repo:   repo,
		bcrypt: bcrypt,
	}
}

// Params: struct auxiliar para no pasar 4 argumentos sueltos
type UpdateUserParams struct {
	ID       int64
	Name     string
	Phone    string
	Password string
}

func (uc *UpdateUser) Execute(params UpdateUserParams) (*entities.User, error) {
	// 1. Obtener el usuario actual de la BD
	// Necesitamos esto para no borrar datos que el usuario NO quiso cambiar (ej: si solo manda el teléfono)
	currentUser, err := uc.repo.GetByID(params.ID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// 2. Aplicar cambios SOLO si vienen datos (Lógica de PATCH)
	
	if params.Name != "" {
		currentUser.Name = params.Name
	}
	
	if params.Phone != "" {
		currentUser.Phone = params.Phone
	}

	// 3. Manejo especial de la Contraseña
	if params.Password != "" {
		hashedPassword, err := uc.bcrypt.HashPassword(params.Password)
		if err != nil {
			return nil, fmt.Errorf("error al encriptar nueva contraseña: %w", err)
		}
		currentUser.Password = hashedPassword
	}

	// 4. Guardar cambios en BD
	err = uc.repo.Update(currentUser)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	// 5. Limpiar password antes de devolver
	currentUser.Password = ""
	
	return currentUser, nil
}