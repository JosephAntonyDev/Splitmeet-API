package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupGroupRoutes(
	r *gin.Engine,
	createGroupCtrl *controllers.CreateGroupController,
	getGroupCtrl *controllers.GetGroupController,
	getMyGroupsCtrl *controllers.GetMyGroupsController,
	updateGroupCtrl *controllers.UpdateGroupController,
	deleteGroupCtrl *controllers.DeleteGroupController,
	inviteMemberCtrl *controllers.InviteMemberController,
	respondInvitationCtrl *controllers.RespondInvitationController,
	getMembersCtrl *controllers.GetMembersController,
	removeMemberCtrl *controllers.RemoveMemberController,
	getPendingInvitationsCtrl *controllers.GetPendingInvitationsController,
	transferOwnershipCtrl *controllers.TransferOwnershipController,
	setMemberRoleCtrl *controllers.SetMemberRoleController,
	joinGroupCtrl *controllers.JoinGroupController,
	jwtSecret string,
) {
	g := r.Group("groups")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Group CRUD
		g.POST("", createGroupCtrl.Handle)
		g.GET("", getMyGroupsCtrl.Handle)
		g.GET("/:id", getGroupCtrl.Handle)
		g.PATCH("/:id", updateGroupCtrl.Handle)
		g.DELETE("/:id", deleteGroupCtrl.Handle)

		// Ownership & Roles
		g.PATCH("/:id/transfer", transferOwnershipCtrl.Handle)
		g.PATCH("/:id/members/:userId/role", setMemberRoleCtrl.Handle)

		// Members
		g.POST("/:id/join", joinGroupCtrl.Handle)
		g.GET("/:id/members", getMembersCtrl.Handle)
		g.POST("/:id/invite", inviteMemberCtrl.Handle)
		g.PATCH("/:id/invitation", respondInvitationCtrl.Handle)
		g.DELETE("/:id/members/:userId", removeMemberCtrl.Handle)

		// Invitations
		g.GET("/invitations/pending", getPendingInvitationsCtrl.Handle)
	}
}
