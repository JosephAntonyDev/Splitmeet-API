package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/infra/controllers"
	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(
	r *gin.Engine,
	createPaymentCtrl *controllers.CreatePaymentController,
	getPaymentCtrl *controllers.GetPaymentController,
	getPaymentsByOutingCtrl *controllers.GetPaymentsByOutingController,
	confirmPaymentCtrl *controllers.ConfirmPaymentController,
	confirmParticipantPaymentCtrl *controllers.ConfirmParticipantPaymentController,
	deletePaymentCtrl *controllers.DeletePaymentController,
	getPaymentSummaryCtrl *controllers.GetPaymentSummaryController,
	jwtSecret string,
) {
	payments := r.Group("payments")
	payments.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Payment CRUD
		payments.POST("", createPaymentCtrl.Handle)
		payments.GET("/outing/:outingId", getPaymentsByOutingCtrl.Handle)
		payments.GET("/outing/:outingId/summary", getPaymentSummaryCtrl.Handle)
		payments.GET("/:id", getPaymentCtrl.Handle)
		payments.PATCH("/:id/confirm", confirmPaymentCtrl.Handle)
		payments.PATCH("/outings/:outing_id/participants/:participant_id/confirm", confirmParticipantPaymentCtrl.Handle)
		payments.DELETE("/:id", deletePaymentCtrl.Handle)
	}
}
