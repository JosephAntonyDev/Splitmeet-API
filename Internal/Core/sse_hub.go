package core

import (
	"encoding/json"
	"fmt"
	"sync"
)

// SSEHub manages Server-Sent Events connections per user
type SSEHub struct {
	mu      sync.RWMutex
	clients map[int64]map[chan string]bool // userID -> set of channels
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[int64]map[chan string]bool),
	}
}

func (h *SSEHub) Register(userID int64, ch chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[chan string]bool)
	}
	h.clients[userID][ch] = true
}

func (h *SSEHub) Unregister(userID int64, ch chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if channels, ok := h.clients[userID]; ok {
		delete(channels, ch)
		if len(channels) == 0 {
			delete(h.clients, userID)
		}
	}
}

// SendToUser sends an SSE event to all connections of a specific user
func (h *SSEHub) SendToUser(userID int64, eventType string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	channels, ok := h.clients[userID]
	if !ok {
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, string(jsonData))

	for ch := range channels {
		select {
		case ch <- message:
		default:
			// Channel is full, skip
		}
	}
}
