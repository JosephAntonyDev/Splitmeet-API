package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	groupRepository "github.com/JosephAntonyDev/splitmeet-api/internal/group/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra/routes"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/repository"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL, notifSvc *core.NotificationService) {
	// Repositories
	outingRepo := repository.NewOutingPostgresql(dbPool.DB)
	groupRepo := groupRepository.NewGroupPostgreSQLRepository(dbPool)
	userRepo := userRepository.NewUserPostgreSQLRepository(dbPool)

	// Use Cases - Outings
	createOutingUC := app.NewCreateOutingUseCase(outingRepo, groupRepo, userRepo, notifSvc)
	getOutingUC := app.NewGetOutingUseCase(outingRepo)
	getOutingsByGroupUC := app.NewGetOutingsByGroupUseCase(outingRepo)
	getMyOutingsUC := app.NewGetMyOutingsUseCase(outingRepo)
	updateOutingUC := app.NewUpdateOutingUseCase(outingRepo)
	deleteOutingUC := app.NewDeleteOutingUseCase(outingRepo)

	// Use Cases - Participants
	addParticipantUC := app.NewAddParticipantUseCase(outingRepo, userRepo, notifSvc)
	getParticipantsUC := app.NewGetParticipantsUseCase(outingRepo)
	confirmParticipationUC := app.NewConfirmParticipationUseCase(outingRepo, userRepo, notifSvc)
	removeParticipantUC := app.NewRemoveParticipantUseCase(outingRepo)
	joinOutingUC := app.NewJoinOutingUseCase(outingRepo)

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
	joinOutingCtrl := controllers.NewJoinOutingController(joinOutingUC)

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
		joinOutingCtrl,
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
