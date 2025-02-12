package services

import (
	"log"

	"github.com/gorilla/websocket"
)

// BroadcastService 负责消息广播
type BroadcastService struct {
	connectionService   *ConnectionService
	subscriptionService *SubscriptionService
}

// NewBroadcastService 创建一个新的广播管理服务
func NewBroadcastService(connectionService *ConnectionService, subscriptionService *SubscriptionService) *BroadcastService {
	return &BroadcastService{
		connectionService:   connectionService,
		subscriptionService: subscriptionService,
	}
}

// Broadcast 广播消息给所有订阅了该消息类型的客户端
func (bs *BroadcastService) Broadcast(messageType string, message []byte) {
	connections := bs.connectionService.GetConnections()

	log.Printf("Broadcasting message of type %s to %d clients", messageType, len(connections))
	for conn := range connections {
		subscriptions := bs.subscriptionService.GetSubscriptions(conn)
		if subscriptions[messageType] {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error broadcasting message:", err)
				bs.connectionService.RemoveConnection(conn)
			}
		}
	}
}
