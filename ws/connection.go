package services

import (
	"log"

	"github.com/gorilla/websocket"
)

// ConnectionService 负责连接管理
type ConnectionService struct {
	clients map[*websocket.Conn]bool
}

// NewConnectionService 创建一个新的连接管理服务
func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddConnection 添加连接
func (cs *ConnectionService) AddConnection(conn *websocket.Conn) {
	cs.clients[conn] = true
	log.Printf("New connection added: %s", conn.RemoteAddr())
}

// RemoveConnection 移除连接
func (cs *ConnectionService) RemoveConnection(conn *websocket.Conn) {
	delete(cs.clients, conn)
	conn.Close()
	log.Printf("Connection removed: %s", conn.RemoteAddr())
}

// GetConnections 获取所有连接
func (cs *ConnectionService) GetConnections() map[*websocket.Conn]bool {
	return cs.clients
}
