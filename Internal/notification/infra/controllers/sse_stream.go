package controllers

import (
	"io"
	"net/http"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/gin-gonic/gin"
)

type SSEStreamController struct {
	hub *core.SSEHub
}

func NewSSEStreamController(hub *core.SSEHub) *SSEStreamController {
	return &SSEStreamController{hub: hub}
}

func (ctrl *SSEStreamController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	uid := userID.(int64)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch := make(chan string, 64)
	ctrl.hub.Register(uid, ch)
	defer func() {
		ctrl.hub.Unregister(uid, ch)
		close(ch)
	}()

	// Send initial connected event
	c.SSEvent("connected", gin.H{"message": "Conexión SSE establecida"})
	c.Writer.Flush()

	clientGone := c.Request.Context().Done()

	for {
		select {
		case <-clientGone:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			io.WriteString(c.Writer, msg)
			c.Writer.Flush()
		}
	}
}
