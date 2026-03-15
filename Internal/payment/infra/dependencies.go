package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/infra/routes"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Repository
	paymentRepo := repository.NewPaymentPostgresql(dbPool.DB)

	// Use Cases
	createPaymentUC := app.NewCreatePaymentUseCase(paymentRepo)
	getPaymentUC := app.NewGetPaymentUseCase(paymentRepo)
	getPaymentsByOutingUC := app.NewGetPaymentsByOutingUseCase(paymentRepo)
	confirmPaymentUC := app.NewConfirmPaymentUseCase(paymentRepo)
	confirmParticipantPaymentUC := app.NewConfirmParticipantPaymentUseCase(paymentRepo)
	deletePaymentUC := app.NewDeletePaymentUseCase(paymentRepo)
	getPaymentSummaryUC := app.NewGetPaymentSummaryUseCase(paymentRepo)

	// Controllers
	createPaymentCtrl := controllers.NewCreatePaymentController(createPaymentUC)
	getPaymentCtrl := controllers.NewGetPaymentController(getPaymentUC)
	getPaymentsByOutingCtrl := controllers.NewGetPaymentsByOutingController(getPaymentsByOutingUC)
	confirmPaymentCtrl := controllers.NewConfirmPaymentController(confirmPaymentUC)
	confirmParticipantPaymentCtrl := controllers.NewConfirmParticipantPaymentController(confirmParticipantPaymentUC)
	deletePaymentCtrl := controllers.NewDeletePaymentController(deletePaymentUC)
	getPaymentSummaryCtrl := controllers.NewGetPaymentSummaryController(getPaymentSummaryUC)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Routes
	routes.SetupPaymentRoutes(
		r,
		createPaymentCtrl,
		getPaymentCtrl,
		getPaymentsByOutingCtrl,
		confirmPaymentCtrl,
		confirmParticipantPaymentCtrl,
		deletePaymentCtrl,
		getPaymentSummaryCtrl,
		jwtSecret,
	)
}
