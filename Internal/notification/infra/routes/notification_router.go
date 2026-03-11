package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/infra/controllers"
	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(
	r *gin.Engine,
	getNotificationsCtrl *controllers.GetNotificationsController,
	markAsReadCtrl *controllers.MarkAsReadController,
	sseStreamCtrl *controllers.SSEStreamController,
	jwtSecret string,
) {
	g := r.Group("notifications")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	{
		g.GET("", getNotificationsCtrl.Handle)
		g.GET("/stream", sseStreamCtrl.Handle)
		g.PATCH("/:id/read", markAsReadCtrl.HandleOne)
		g.PATCH("/read-all", markAsReadCtrl.HandleAll)
	}
}
