package websocket

import (
	"net/http"
	"os"
	"strings"
	"time"

	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 从环境变量获取允许的来源，生产环境应配置具体域名
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			// 开发环境允许所有来源
			if os.Getenv("GO_ENV") != "production" && os.Getenv("GIN_MODE") != "release" {
				return true
			}
			// 生产环境默认只允许同源
			return r.Header.Get("Origin") == ""
		}
		origin := r.Header.Get("Origin")
		// 支持多个来源，用逗号分隔
		for _, allowed := range strings.Split(allowedOrigin, ",") {
			if strings.TrimSpace(allowed) == origin {
				return true
			}
		}
		return false
	},
}

const (
	// 写入超时
	writeWait = 10 * time.Second
	// 心跳超时
	pongWait = 60 * time.Second
	// 心跳间隔
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小
	maxMessageSize = 8192
)

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	// 从 token 获取用户信息
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := jwttoken.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Errorf("Failed to upgrade WebSocket: %v", err)
		return
	}

	hub := GetHub()
	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   0, // 从数据库获取用户ID
		Username: claims.Username,
	}

	hub.Register(client)

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// readPump 从 WebSocket 读取消息
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Errorf("WebSocket read error: %v", err)
			}
			break
		}
		// 处理收到的消息（如心跳、命令等）
		c.handleMessage(message)
	}
}

// writePump 向 WebSocket 写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了 channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 合并队列中的消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(message []byte) {
	// 可以在这里处理客户端发来的消息
	// 例如：标记通知已读、确认收到等
	slog.Debugf("Received message from user %s: %s", c.Username, string(message))
}
