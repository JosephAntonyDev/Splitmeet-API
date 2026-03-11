package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/notification/app"
	"github.com/gin-gonic/gin"
)

type GetNotificationsController struct {
	useCase *app.GetNotifications
}

func NewGetNotificationsController(useCase *app.GetNotifications) *GetNotificationsController {
	return &GetNotificationsController{useCase: useCase}
}

func (ctrl *GetNotificationsController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	params := core.GetPaginationParams(c)

	input := app.GetNotificationsInput{
		UserID: userID.(int64),
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	notifications, total, err := ctrl.useCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}

	var data []gin.H
	for _, n := range notifications {
		item := gin.H{
			"id":              n.ID,
			"type":            n.Type,
			"title":           n.Title,
			"message":         n.Message,
			"inviter_name":    n.InviterName,
			"group_name":      n.GroupName,
			"outing_name":     n.OutingName,
			"is_read":         n.IsRead,
			"response_status": n.ResponseStatus,
			"created_at":      n.CreatedAt,
		}
		if n.ReferenceID != nil {
			item["reference_id"] = *n.ReferenceID
		}
		data = append(data, item)
	}

	if data == nil {
		data = []gin.H{}
	}

	c.JSON(http.StatusOK, core.NewPaginatedResponse(data, params.Page, params.Limit, total))
}
