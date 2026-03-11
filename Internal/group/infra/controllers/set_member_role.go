package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type SetMemberRoleController struct {
	useCase *app.SetMemberRole
}

func NewSetMemberRoleController(useCase *app.SetMemberRole) *SetMemberRoleController {
	return &SetMemberRoleController{useCase: useCase}
}

type SetMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (ctrl *SetMemberRoleController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de grupo inválido"})
		return
	}

	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	var req SetMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El rol es requerido"})
		return
	}

	input := app.SetMemberRoleInput{
		GroupID:      groupID,
		TargetUserID: targetUserID,
		Role:         req.Role,
		RequestedBy:  userID.(int64),
	}

	err = ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rol actualizado exitosamente"})
}
