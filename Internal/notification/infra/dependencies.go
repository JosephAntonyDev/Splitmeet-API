package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/infra/routes"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL, hub *core.SSEHub) {
	// Repository
	notifRepo := repository.NewNotificationPostgreSQLRepository(dbPool)

	// Use Cases
	getNotificationsUC := app.NewGetNotifications(notifRepo)
	markAsReadUC := app.NewMarkAsRead(notifRepo)
	markAllAsReadUC := app.NewMarkAllAsRead(notifRepo)
	registerDeviceTokenUC := app.NewRegisterDeviceToken(notifRepo)

	// Controllers
	getNotificationsCtrl := controllers.NewGetNotificationsController(getNotificationsUC)
	markAsReadCtrl := controllers.NewMarkAsReadController(markAsReadUC, markAllAsReadUC)
	sseStreamCtrl := controllers.NewSSEStreamController(hub)
	registerDeviceTokenCtrl := controllers.NewRegisterDeviceTokenController(registerDeviceTokenUC)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Routes
	routes.SetupNotificationRoutes(r, getNotificationsCtrl, markAsReadCtrl, sseStreamCtrl, registerDeviceTokenCtrl, jwtSecret)
}
