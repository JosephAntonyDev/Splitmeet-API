package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type RespondInvitationController struct {
	useCase *app.RespondInvitation
}

func NewRespondInvitationController(useCase *app.RespondInvitation) *RespondInvitationController {
	return &RespondInvitationController{useCase: useCase}
}

type RespondInvitationRequest struct {
	Accept bool `json:"accept"`
}

func (ctrl *RespondInvitationController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	idParam := c.Param("id")
	groupID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de grupo inválido"})
		return
	}

	var req RespondInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	input := app.RespondInvitationInput{
		GroupID: groupID,
		UserID:  userID.(int64),
		Accept:  req.Accept,
	}

	err = ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var message string
	if req.Accept {
		message = "Te has unido al grupo exitosamente"
	} else {
		message = "Invitación rechazada"
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
