package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket 升级配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境中建议设置具体的域名
	},
}

// 消息类型常量
const (
	ChatMessageType   = "chat"
	SystemMessageType = "system"
	NotificationType  = "notification"
)

// 订阅和取消订阅消息结构
type SubscriptionMessage struct {
	Action  string `json:"action"`  // 操作类型，"subscribe" 或 "unsubscribe"
	Content string `json:"content"` // 消息类型（如 "chat", "system", "notification"）
	Status  string `json:"status"`  // 消息状态，例："SENT"
}

// 聊天消息结构
type ChatMessage struct {
	Type     string `json:"type"`     // 消息类型：'chat', 'system', 'notification'
	User     string `json:"user"`     // 发送者
	Message  string `json:"message"`  // 消息内容
	SendTime int64  `json:"sendTime"` // 发送时间戳
	Status   string `json:"status"`   // 消息状态，例："SENT"
}

// WebSocket 连接管理
type WebSocketManager struct {
	clients map[*websocket.Conn]map[string]bool // 存储每个连接的客户端及其订阅的消息类型
	mutex   sync.Mutex                          // 保护 clients 访问的互斥锁
}

// 创建 WebSocket 管理器
var wsManager = WebSocketManager{
	clients: make(map[*websocket.Conn]map[string]bool),
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接到 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 添加连接到管理器
	if !wsManager.addClient(conn) {
		return // 如果连接已满，直接返回
	}
	defer wsManager.removeClient(conn)

	log.Println("New WebSocket connection established.")

	for {
		// 读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break // 退出循环，关闭连接
		}

		log.Printf("Received message: %s", message)

		// 解析消息，识别操作类型并执行相应操作
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// 判断是订阅/取消订阅消息还是聊天消息
		if _, exists := msg["action"]; exists {
			var subMsg SubscriptionMessage
			if err := json.Unmarshal(message, &subMsg); err == nil {
				switch subMsg.Action {
				case "subscribe":
					// 处理订阅消息
					wsManager.subscribeToMessageType(conn, subMsg.Content)
				case "unsubscribe":
					// 处理取消订阅消息
					wsManager.unsubscribeFromMessageType(conn, subMsg.Content)
				}
			}
		} else {
			// 如果没有 "action" 字段，则处理为聊天消息
			var chatMsg ChatMessage
			if err := json.Unmarshal(message, &chatMsg); err == nil {
				// 广播聊天消息给订阅的客户端
				wsManager.broadcast(chatMsg.Type, message)
			} else {
				log.Println("Invalid message format")
			}
		}
	}
}

// addClient 将新的 WebSocket 连接添加到管理器
const MaxClients = 10000

func (wm *WebSocketManager) addClient(conn *websocket.Conn) bool {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if len(wm.clients) >= MaxClients {
		log.Println("Max WebSocket clients reached, rejecting new connection.")
		conn.Close()
		return false
	}

	// 初始化客户端的订阅信息
	wm.clients[conn] = make(map[string]bool)
	log.Println("Client added. Total clients:", len(wm.clients))
	return true
}

// removeClient 关闭 WebSocket 连接并从管理器移除
func (wm *WebSocketManager) removeClient(conn *websocket.Conn) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	delete(wm.clients, conn)
	conn.Close()
	log.Println("Client removed. Total clients:", len(wm.clients))
}

// subscribeToMessageType 让客户端订阅指定类型的消息
func (wm *WebSocketManager) subscribeToMessageType(conn *websocket.Conn, messageType string) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if _, exists := wm.clients[conn]; exists {
		wm.clients[conn][messageType] = true
		log.Printf("Client %v subscribed to %s messages", conn.RemoteAddr(), messageType)
	} else {
		log.Println("Client not found when subscribing")
	}
}

// unsubscribeFromMessageType 让客户端取消订阅指定类型的消息
func (wm *WebSocketManager) unsubscribeFromMessageType(conn *websocket.Conn, messageType string) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if _, exists := wm.clients[conn]; exists {
		delete(wm.clients[conn], messageType)
		log.Printf("Client %v unsubscribed from %s messages", conn.RemoteAddr(), messageType)
	} else {
		log.Println("Client not found when unsubscribing")
	}
}

// broadcast 向所有连接的客户端发送指定类型的消息
func (wm *WebSocketManager) broadcast(messageType string, message []byte) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	var broadcastCount int
	for client, subscriptions := range wm.clients {
		// 只发送消息给订阅了该消息类型的客户端
		if subscriptions[messageType] {
			log.Printf("Sending message to client %v with subscriptions: %v", client.RemoteAddr(), subscriptions)
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error broadcasting message:", err)
				client.Close()
				delete(wm.clients, client)
			} else {
				broadcastCount++
			}
		}
	}

	log.Printf("Broadcasted %s message to %d clients\n", messageType, broadcastCount)
}
