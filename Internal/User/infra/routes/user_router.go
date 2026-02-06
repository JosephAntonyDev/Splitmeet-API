package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, createUserCtrl *controllers.CreateUserController, loginUserCtrl *controllers.LoginUserController,
	getUserCtrl *controllers.GetUserController, getByUsernameCtrl *controllers.GetByUsernameController, getProfileCtrl *controllers.GetProfileController, updateUserCtrl *controllers.UpdateUserController,
	deleteUserCtrl *controllers.DeleteUserController, searchUsersCtrl *controllers.SearchUsersController,
	getPendingInvitationsCtrl *controllers.GetPendingInvitationsController,
	jwtSecret string) {
	g := r.Group("users")
	{
		g.POST("", createUserCtrl.Handle)
		g.POST("/login", loginUserCtrl.Handle)
	}
	gPrivate := r.Group("users")
	gPrivate.Use(middleware.AuthMiddleware(jwtSecret))
	{
		gPrivate.GET("/get/:id", getUserCtrl.Handle)
		gPrivate.GET("/username/:username", getByUsernameCtrl.Handle)
		gPrivate.GET("/search", searchUsersCtrl.Handle)
		gPrivate.GET("/profile", getProfileCtrl.Handle)
		gPrivate.GET("/invitations", getPendingInvitationsCtrl.Handle)
		gPrivate.PATCH("/update", updateUserCtrl.Handle)
		gPrivate.DELETE("/delete", deleteUserCtrl.Handle)
	}
}
