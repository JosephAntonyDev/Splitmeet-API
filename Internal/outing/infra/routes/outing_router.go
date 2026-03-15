package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/controllers"
	"github.com/gin-gonic/gin"
)

func SetupOutingRoutes(
	r *gin.Engine,
	createOutingCtrl *controllers.CreateOutingController,
	getOutingCtrl *controllers.GetOutingController,
	getOutingsByGroupCtrl *controllers.GetOutingsByGroupController,
	getMyOutingsCtrl *controllers.GetMyOutingsController,
	updateOutingCtrl *controllers.UpdateOutingController,
	deleteOutingCtrl *controllers.DeleteOutingController,
	addParticipantCtrl *controllers.AddParticipantController,
	getParticipantsCtrl *controllers.GetParticipantsController,
	confirmParticipationCtrl *controllers.ConfirmParticipationController,
	removeParticipantCtrl *controllers.RemoveParticipantController,
	joinOutingCtrl *controllers.JoinOutingController,
	addItemCtrl *controllers.AddItemController,
	getItemsCtrl *controllers.GetItemsController,
	updateItemCtrl *controllers.UpdateItemController,
	removeItemCtrl *controllers.RemoveItemController,
	setItemSplitsCtrl *controllers.SetItemSplitsController,
	getItemSplitsCtrl *controllers.GetItemSplitsController,
	calculateSplitsCtrl *controllers.CalculateSplitsController,
	jwtSecret string,
) {
	outings := r.Group("outings")
	outings.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Outing CRUD
		outings.POST("", createOutingCtrl.Handle)
		outings.GET("/me", getMyOutingsCtrl.Handle)
		outings.GET("/group/:groupId", getOutingsByGroupCtrl.Handle)
		outings.GET("/:id", getOutingCtrl.Handle)
		outings.PATCH("/:id", updateOutingCtrl.Handle)
		outings.DELETE("/:id", deleteOutingCtrl.Handle)

		// Participants
		outings.POST("/:id/join", joinOutingCtrl.Handle)
		outings.POST("/:id/participants", addParticipantCtrl.Handle)
		outings.GET("/:id/participants", getParticipantsCtrl.Handle)
		outings.PATCH("/:id/participants/confirm", confirmParticipationCtrl.Handle)
		outings.DELETE("/:id/participants/:userId", removeParticipantCtrl.Handle)

		// Items
		outings.POST("/:id/items", addItemCtrl.Handle)
		outings.GET("/:id/items", getItemsCtrl.Handle)
		outings.PATCH("/:id/items/:itemId", updateItemCtrl.Handle)
		outings.DELETE("/:id/items/:itemId", removeItemCtrl.Handle)

		// Splits
		outings.POST("/:id/items/:itemId/splits", setItemSplitsCtrl.Handle)
		outings.GET("/:id/items/:itemId/splits", getItemSplitsCtrl.Handle)

		// Calculate totals
		outings.GET("/:id/calculate", calculateSplitsCtrl.Handle)
	}
}
