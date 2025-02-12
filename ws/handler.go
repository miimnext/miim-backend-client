package services

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// HandleMessage 处理 WebSocket 消息
func HandleMessage(conn *websocket.Conn, message []byte, subscriptionService *SubscriptionService) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}

	// 判断是否是订阅或取消订阅操作
	if action, exists := msg["action"]; exists {
		switch action {
		case "subscribe":
			messageType := msg["content"].(string)
			subscriptionService.Subscribe(conn, messageType)
		case "unsubscribe":
			messageType := msg["content"].(string)
			subscriptionService.Unsubscribe(conn, messageType)
		}
	} else {
		// 普通消息处理（聊天）
		log.Printf("Received message: %s", string(message))
		// 可以扩展：根据消息类型进行不同处理
	}
}
