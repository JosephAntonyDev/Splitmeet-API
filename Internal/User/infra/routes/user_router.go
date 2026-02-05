package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
)

func SetupUserRoutes(r *gin.Engine, createUserCtrl *controllers.CreateUserController, jwtSecret string) {
	g := r.Group("users")
	{
		g.POST("", createUserCtrl.Handle)
	}
	gPrivate := r.Group("users")
	gPrivate.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Rutas privadas aquí
	}
}