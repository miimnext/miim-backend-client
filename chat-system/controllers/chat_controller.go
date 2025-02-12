package controllers

import (
	"chat-system/config"
	"chat-system/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader，用于将 HTTP 请求升级为 WebSocket 协议
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket 连接池，用于存储所有的用户连接
var wsConnections = make(map[string]*websocket.Conn)

// WebSocketHandler 处理 WebSocket 连接
func WebSocketHandler(c *gin.Context) {
	// 获取用户ID，通常通过认证获取
	userID := c.Param("user_id")

	// 升级 HTTP 请求为 WebSocket 协议
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// 将连接存入连接池
	wsConnections[userID] = conn
	log.Printf("User %s connected", userID)

	// 不断读取客户端发送的消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(wsConnections, userID)
			break
		}

		// 处理接收到的消息
		err = handleMessage(userID, messageType, message)
		if err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

// 处理接收到的消息
func handleMessage(senderID string, messageType int, message []byte) error {
	// 根据消息类型、内容进行处理，存入数据库
	// 这里可以根据业务需求进一步扩展处理逻辑

	// 例如：处理文本消息
	messageContent := string(message)
	msg := models.Message{
		MessageID:      fmt.Sprintf("%s-%d", senderID, time.Now().UnixNano()), // 可以根据需求生成唯一的MessageID
		ConversationID: "",                                                    // 如果是群聊，填充群ID；如果是单聊，填充用户对话ID
		SenderID:       senderID,
		ReceiverID:     "", // 可以指定接收者的ID
		Content:        messageContent,
		MessageType:    "text",
		Status:         "sent", // 默认消息状态为发送
	}

	// 将消息存入数据库
	if err := config.DB.Create(&msg).Error; err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}

	// 广播消息到所有 WebSocket 连接
	return broadcastMessage(senderID, messageType, message)
}

// 广播消息到所有 WebSocket 连接
func broadcastMessage(senderID string, messageType int, message []byte) error {
	for userID, conn := range wsConnections {
		if userID != senderID {
			err := conn.WriteMessage(messageType, message)
			if err != nil {
				log.Printf("Error broadcasting to %s: %v", userID, err)
			}
		}
	}
	return nil
}
