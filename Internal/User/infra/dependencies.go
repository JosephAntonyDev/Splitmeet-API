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
	loginUserUseCase := app.NewLoginUser(userRepo, bcryptService, adapters.NewJWTManager(os.Getenv("JWT_SECRET")))
	getUserUseCase := app.NewGetUser(userRepo)
	getProfileUseCase := app.NewGetProfile(userRepo)
	updateUserUseCase := app.NewUpdateUser(userRepo, bcryptService)
	deleteUserUseCase := app.NewDeleteUser(userRepo)
	

	createUserController := controllers.NewCreateUserController(createUserUseCase)
	loginUserController := controllers.NewLoginUserController(loginUserUseCase)
	getUserController := controllers.NewGetUserController(getUserUseCase)
	getProfileController := controllers.NewGetProfileController(getProfileUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)

	jwtSecret := os.Getenv("JWT_SECRET")

	routes.SetupUserRoutes(r, createUserController, loginUserController, getUserController, getProfileController, updateUserController, deleteUserController, jwtSecret)
}