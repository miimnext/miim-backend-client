package models

import (
	"gorm.io/gorm"
)

// WSConnection WebSocket 连接模型
// WebSocket 连接表存储每个连接的信息
type WSConnection struct {
	gorm.Model
	ConnectionID string `json:"connection_id" gorm:"primaryKey"` // 连接ID，主键
	UserID       string `json:"user_id"`                         // 用户ID
	ConnectedAt  string `json:"connected_at"`                    // 连接时间
}
