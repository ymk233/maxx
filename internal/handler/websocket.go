package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/Bowl42/maxx-next/internal/domain"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要严格检查
	},
}

type WSMessage struct {
	Type string      `json:"type"` // "proxy_request_update", "stats_update"
	Data interface{} `json:"data"`
}

type WebSocketHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan WSMessage
	mu        sync.RWMutex
}

func NewWebSocketHub() *WebSocketHub {
	hub := &WebSocketHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan WSMessage, 100),
	}
	go hub.run()
	return hub
}

func (h *WebSocketHub) run() {
	for msg := range h.broadcast {
		h.mu.RLock()
		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.RUnlock()
	}
}

func (h *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	// 保持连接，处理客户端消息（心跳等）
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *WebSocketHub) BroadcastProxyRequest(req *domain.ProxyRequest) {
	h.broadcast <- WSMessage{
		Type: "proxy_request_update",
		Data: req,
	}
}

func (h *WebSocketHub) BroadcastStats(stats interface{}) {
	h.broadcast <- WSMessage{
		Type: "stats_update",
		Data: stats,
	}
}
