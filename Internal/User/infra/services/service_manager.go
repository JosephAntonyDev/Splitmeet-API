package services

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/ports"
	adapters "github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/adapters"
)

// Inicializar el servicio de BCrypt
func InitBcryptService() ports.IBcryptService {
	return adapters.NewBcrypt()
}

// Inicializar el Token Manager
func InitTokenManager() ports.TokenManager {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET no está configurado en las variables de entorno")
	}
	return &adapters.JWTManager{SecretKey: jwtSecret}
}