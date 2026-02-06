package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/gin-gonic/gin"

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
	getByUsernameUseCase := app.NewGetByUsername(userRepo)
	getProfileUseCase := app.NewGetProfile(userRepo)
	updateUserUseCase := app.NewUpdateUser(userRepo, bcryptService)
	deleteUserUseCase := app.NewDeleteUser(userRepo)
	searchUsersUseCase := app.NewSearchUsers(userRepo)
	getPendingInvitationsUseCase := app.NewGetPendingInvitations(dbPool.DB)

	createUserController := controllers.NewCreateUserController(createUserUseCase)
	loginUserController := controllers.NewLoginUserController(loginUserUseCase)
	getUserController := controllers.NewGetUserController(getUserUseCase)
	getByUsernameController := controllers.NewGetByUsernameController(getByUsernameUseCase)
	getProfileController := controllers.NewGetProfileController(getProfileUseCase)
	updateUserController := controllers.NewUpdateUserController(updateUserUseCase)
	deleteUserController := controllers.NewDeleteUserController(deleteUserUseCase)
	searchUsersController := controllers.NewSearchUsersController(searchUsersUseCase)
	getPendingInvitationsController := controllers.NewGetPendingInvitationsController(getPendingInvitationsUseCase)

	jwtSecret := os.Getenv("JWT_SECRET")

	routes.SetupUserRoutes(r, createUserController, loginUserController, getUserController, getByUsernameController, getProfileController, updateUserController, deleteUserController, searchUsersController, getPendingInvitationsController, jwtSecret)
}
