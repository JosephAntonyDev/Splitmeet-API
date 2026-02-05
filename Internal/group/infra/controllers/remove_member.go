package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type RemoveMemberController struct {
	useCase *app.RemoveMember
}

func NewRemoveMemberController(useCase *app.RemoveMember) *RemoveMemberController {
	return &RemoveMemberController{useCase: useCase}
}

func (ctrl *RemoveMemberController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	groupIDParam := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de grupo inválido"})
		return
	}

	memberIDParam := c.Param("userId")
	memberID, err := strconv.ParseInt(memberIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	input := app.RemoveMemberInput{
		GroupID:        groupID,
		MemberToRemove: memberID,
		RequestedBy:    userID.(int64),
	}

	err = ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Miembro removido exitosamente"})
}
