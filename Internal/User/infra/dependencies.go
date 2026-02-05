package infra

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/adapters"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/routes"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	
	userRepo := repository.NewUserPostgreSQLRepository(dbPool)

	bcryptService := adapters.NewBcrypt()

	createUserUseCase := app.NewCreateUser(userRepo, bcryptService)

	createUserController := controllers.NewCreateUserController(createUserUseCase)

	jwtSecret := os.Getenv("JWT_SECRET")

	routes.SetupUserRoutes(r, createUserController, jwtSecret)
}