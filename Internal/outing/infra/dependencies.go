package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/routes"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Repository
	outingRepo := repository.NewOutingPostgresql(dbPool.DB)

	// Use Cases - Outings
	createOutingUC := app.NewCreateOutingUseCase(outingRepo)
	getOutingUC := app.NewGetOutingUseCase(outingRepo)
	getOutingsByGroupUC := app.NewGetOutingsByGroupUseCase(outingRepo)
	getMyOutingsUC := app.NewGetMyOutingsUseCase(outingRepo)
	updateOutingUC := app.NewUpdateOutingUseCase(outingRepo)
	deleteOutingUC := app.NewDeleteOutingUseCase(outingRepo)

	// Use Cases - Participants
	addParticipantUC := app.NewAddParticipantUseCase(outingRepo)
	getParticipantsUC := app.NewGetParticipantsUseCase(outingRepo)
	confirmParticipationUC := app.NewConfirmParticipationUseCase(outingRepo)
	removeParticipantUC := app.NewRemoveParticipantUseCase(outingRepo)

	// Use Cases - Items
	addItemUC := app.NewAddItemUseCase(outingRepo)
	getItemsUC := app.NewGetItemsUseCase(outingRepo)
	updateItemUC := app.NewUpdateItemUseCase(outingRepo)
	removeItemUC := app.NewRemoveItemUseCase(outingRepo)

	// Use Cases - Splits
	setItemSplitsUC := app.NewSetItemSplitsUseCase(outingRepo)
	getItemSplitsUC := app.NewGetItemSplitsUseCase(outingRepo)
	calculateSplitsUC := app.NewCalculateSplitsUseCase(outingRepo)

	// Controllers - Outings
	createOutingCtrl := controllers.NewCreateOutingController(createOutingUC)
	getOutingCtrl := controllers.NewGetOutingController(getOutingUC)
	getOutingsByGroupCtrl := controllers.NewGetOutingsByGroupController(getOutingsByGroupUC)
	getMyOutingsCtrl := controllers.NewGetMyOutingsController(getMyOutingsUC)
	updateOutingCtrl := controllers.NewUpdateOutingController(updateOutingUC)
	deleteOutingCtrl := controllers.NewDeleteOutingController(deleteOutingUC)

	// Controllers - Participants
	addParticipantCtrl := controllers.NewAddParticipantController(addParticipantUC)
	getParticipantsCtrl := controllers.NewGetParticipantsController(getParticipantsUC)
	confirmParticipationCtrl := controllers.NewConfirmParticipationController(confirmParticipationUC)
	removeParticipantCtrl := controllers.NewRemoveParticipantController(removeParticipantUC)

	// Controllers - Items
	addItemCtrl := controllers.NewAddItemController(addItemUC)
	getItemsCtrl := controllers.NewGetItemsController(getItemsUC)
	updateItemCtrl := controllers.NewUpdateItemController(updateItemUC)
	removeItemCtrl := controllers.NewRemoveItemController(removeItemUC)

	// Controllers - Splits
	setItemSplitsCtrl := controllers.NewSetItemSplitsController(setItemSplitsUC)
	getItemSplitsCtrl := controllers.NewGetItemSplitsController(getItemSplitsUC)
	calculateSplitsCtrl := controllers.NewCalculateSplitsController(calculateSplitsUC)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Routes
	routes.SetupOutingRoutes(
		r,
		createOutingCtrl,
		getOutingCtrl,
		getOutingsByGroupCtrl,
		getMyOutingsCtrl,
		updateOutingCtrl,
		deleteOutingCtrl,
		addParticipantCtrl,
		getParticipantsCtrl,
		confirmParticipationCtrl,
		removeParticipantCtrl,
		addItemCtrl,
		getItemsCtrl,
		updateItemCtrl,
		removeItemCtrl,
		setItemSplitsCtrl,
		getItemSplitsCtrl,
		calculateSplitsCtrl,
		jwtSecret,
	)
}
