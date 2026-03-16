package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type JoinGroupController struct {
	useCase *app.JoinGroupUseCase
}

func NewJoinGroupController(useCase *app.JoinGroupUseCase) *JoinGroupController {
	return &JoinGroupController{useCase: useCase}
}

func (c *JoinGroupController) Handle(ctx *gin.Context) {
	// Parse group ID from URL
	groupIDStr := ctx.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	// Get authenticated user ID
	userIDRaw, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// Execute use case
	member, err := c.useCase.Execute(groupID, userID)
	if err != nil {
		// Podría ser StatusConflict si ya es miembro, pero Bad Request es suficiente para este caso
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully joined group",
		"data":    member,
	})
}
