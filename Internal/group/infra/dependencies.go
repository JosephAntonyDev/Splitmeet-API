package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/infra/routes"
	userRepository "github.com/JosephAntonyDev/splitmeet-api/internal/user/infra/repository"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL, notifSvc *core.NotificationService) {
	// Repositories
	groupRepo := repository.NewGroupPostgreSQLRepository(dbPool)
	userRepo := userRepository.NewUserPostgreSQLRepository(dbPool)

	// Use Cases
	createGroupUseCase := app.NewCreateGroup(groupRepo)
	getGroupUseCase := app.NewGetGroup(groupRepo)
	getMyGroupsUseCase := app.NewGetMyGroups(groupRepo)
	updateGroupUseCase := app.NewUpdateGroup(groupRepo)
	deleteGroupUseCase := app.NewDeleteGroup(groupRepo)
	inviteMemberUseCase := app.NewInviteMember(groupRepo, userRepo, notifSvc)
	respondInvitationUseCase := app.NewRespondInvitation(groupRepo, userRepo, notifSvc)
	getMembersUseCase := app.NewGetMembers(groupRepo)
	removeMemberUseCase := app.NewRemoveMember(groupRepo)
	getPendingInvitationsUseCase := app.NewGetPendingInvitations(groupRepo)
	transferOwnershipUseCase := app.NewTransferOwnership(groupRepo)
	setMemberRoleUseCase := app.NewSetMemberRole(groupRepo)
	joinGroupUseCase := app.NewJoinGroupUseCase(groupRepo)

	// Controllers
	createGroupController := controllers.NewCreateGroupController(createGroupUseCase)
	getGroupController := controllers.NewGetGroupController(getGroupUseCase)
	getMyGroupsController := controllers.NewGetMyGroupsController(getMyGroupsUseCase)
	updateGroupController := controllers.NewUpdateGroupController(updateGroupUseCase)
	deleteGroupController := controllers.NewDeleteGroupController(deleteGroupUseCase)
	inviteMemberController := controllers.NewInviteMemberController(inviteMemberUseCase)
	respondInvitationController := controllers.NewRespondInvitationController(respondInvitationUseCase)
	getMembersController := controllers.NewGetMembersController(getMembersUseCase)
	removeMemberController := controllers.NewRemoveMemberController(removeMemberUseCase)
	getPendingInvitationsController := controllers.NewGetPendingInvitationsController(getPendingInvitationsUseCase)
	transferOwnershipController := controllers.NewTransferOwnershipController(transferOwnershipUseCase)
	setMemberRoleController := controllers.NewSetMemberRoleController(setMemberRoleUseCase)
	joinGroupController := controllers.NewJoinGroupController(joinGroupUseCase)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Routes
	routes.SetupGroupRoutes(
		r,
		createGroupController,
		getGroupController,
		getMyGroupsController,
		updateGroupController,
		deleteGroupController,
		inviteMemberController,
		respondInvitationController,
		getMembersController,
		removeMemberController,
		getPendingInvitationsController,
		transferOwnershipController,
		setMemberRoleController,
		joinGroupController,
		jwtSecret,
	)
}
