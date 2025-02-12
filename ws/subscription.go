package services

import (
	"sync"

	"github.com/gorilla/websocket"
)

// SubscriptionService 负责客户端的订阅管理
type SubscriptionService struct {
	subscriptions map[*websocket.Conn]map[string]bool
	mutex         sync.Mutex
}

// NewSubscriptionService 创建新的订阅管理服务
func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{
		subscriptions: make(map[*websocket.Conn]map[string]bool),
	}
}

// Subscribe 让客户端订阅指定类型的消息
func (ss *SubscriptionService) Subscribe(conn *websocket.Conn, messageType string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if _, exists := ss.subscriptions[conn]; !exists {
		ss.subscriptions[conn] = make(map[string]bool)
	}

	ss.subscriptions[conn][messageType] = true
}

// Unsubscribe 让客户端取消订阅指定类型的消息
func (ss *SubscriptionService) Unsubscribe(conn *websocket.Conn, messageType string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if _, exists := ss.subscriptions[conn]; exists {
		delete(ss.subscriptions[conn], messageType)
	}
}

// GetSubscriptions 获取客户端的订阅信息
func (ss *SubscriptionService) GetSubscriptions(conn *websocket.Conn) map[string]bool {
	return ss.subscriptions[conn]
}
