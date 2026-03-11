package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/app"
	"github.com/gin-gonic/gin"
)

type MarkAsReadController struct {
	useCase    *app.MarkAsRead
	allUseCase *app.MarkAllAsRead
}

func NewMarkAsReadController(useCase *app.MarkAsRead, allUseCase *app.MarkAllAsRead) *MarkAsReadController {
	return &MarkAsReadController{useCase: useCase, allUseCase: allUseCase}
}

func (ctrl *MarkAsReadController) HandleOne(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificación inválido"})
		return
	}

	err = ctrl.useCase.Execute(notificationID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notificación marcada como leída"})
}

func (ctrl *MarkAsReadController) HandleAll(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	err := ctrl.allUseCase.Execute(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todas las notificaciones marcadas como leídas"})
}
