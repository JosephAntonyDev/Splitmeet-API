package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/group/app"
	"github.com/gin-gonic/gin"
)

type TransferOwnershipController struct {
	useCase *app.TransferOwnership
}

func NewTransferOwnershipController(useCase *app.TransferOwnership) *TransferOwnershipController {
	return &TransferOwnershipController{useCase: useCase}
}

type TransferOwnershipRequest struct {
	NewOwnerID int64 `json:"new_owner_id" binding:"required"`
}

func (ctrl *TransferOwnershipController) Handle(c *gin.Context) {
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

	var req TransferOwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID del nuevo propietario es requerido"})
		return
	}

	input := app.TransferOwnershipInput{
		GroupID:     groupID,
		NewOwnerID:  req.NewOwnerID,
		RequestedBy: userID.(int64),
	}

	err = ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Propiedad transferida exitosamente"})
}
