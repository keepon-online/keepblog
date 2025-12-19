package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
)

// Client WebSocket 客户端
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   int64
	Username string
}

// Hub WebSocket 连接管理中心
type Hub struct {
	// 所有已连接的客户端
	clients map[*Client]bool
	// 按用户ID索引的客户端
	userClients map[int64]*Client
	// 广播频道
	broadcast chan []byte
	// 注册频道
	register chan *Client
	// 注销频道
	unregister chan *Client
	// 互斥锁
	mu sync.RWMutex
}

// Message WebSocket 消息结构
type Message struct {
	Type    string      `json:"type"`    // 消息类型: notice, message, todo, ping, pong
	Payload interface{} `json:"payload"` // 消息内容
	Time    int64       `json:"time"`    // 时间戳
}

var (
	hubInstance *Hub
	hubOnce     sync.Once
)

// GetHub 获取 Hub 单例
func GetHub() *Hub {
	hubOnce.Do(func() {
		hubInstance = &Hub{
			clients:     make(map[*Client]bool),
			userClients: make(map[int64]*Client),
			broadcast:   make(chan []byte),
			register:    make(chan *Client),
			unregister:  make(chan *Client),
		}
	})
	return hubInstance
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.UserID > 0 {
				h.userClients[client.UserID] = client
			}
			h.mu.Unlock()
			slog.Infof("WebSocket client connected: user=%d", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.UserID > 0 {
					delete(h.userClients, client.UserID)
				}
				close(client.Send)
			}
			h.mu.Unlock()
			slog.Infof("WebSocket client disconnected: user=%d", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast 广播消息给所有客户端
func (h *Hub) Broadcast(msg Message) {
	msg.Time = time.Now().Unix()
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Errorf("Failed to marshal broadcast message: %v", err)
		return
	}
	h.broadcast <- data
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, msg Message) bool {
	h.mu.RLock()
	client, ok := h.userClients[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	msg.Time = time.Now().Unix()
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Errorf("Failed to marshal user message: %v", err)
		return false
	}

	select {
	case client.Send <- data:
		return true
	default:
		return false
	}
}

// GetOnlineUsers 获取在线用户ID列表
func (h *Hub) GetOnlineUsers() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]int64, 0, len(h.userClients))
	for userID := range h.userClients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userClients[userID]
	return ok
}
