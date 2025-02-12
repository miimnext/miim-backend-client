package controllers

import (
	services "go_core/ws"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket 升级配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket 处理 WebSocket 请求
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 获取服务实例
	connectionService := services.NewConnectionService()
	subscriptionService := services.NewSubscriptionService()
	broadcastService := services.NewBroadcastService(connectionService, subscriptionService)

	// 添加客户端连接
	connectionService.AddConnection(conn)
	defer connectionService.RemoveConnection(conn)

	// 监听并处理客户端消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// 处理接收到的消息
		services.HandleMessage(conn, message, subscriptionService)

		// 广播消息
		broadcastService.Broadcast("chat", message)
	}
}
