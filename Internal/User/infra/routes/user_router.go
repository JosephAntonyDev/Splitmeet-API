package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
)

func SetupUserRoutes(r *gin.Engine, createUserCtrl *controllers.CreateUserController, loginUserCtrl *controllers.LoginUserController,
	getUserCtrl *controllers.GetUserController, getProfileCtrl *controllers.GetProfileController, updateUserCtrl *controllers.UpdateUserController,
	deleteUserCtrl *controllers.DeleteUserController,
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
		gPrivate.GET("/profile", getProfileCtrl.Handle)
		gPrivate.PATCH("/update", updateUserCtrl.Handle)
		gPrivate.DELETE("/delete", deleteUserCtrl.Handle)
	}
}