package models

import (
	"gorm.io/gorm"
)

// Conversation 会话模型
// 会话表存储用户之间的对话信息
type Conversation struct {
	gorm.Model
	ConversationID string `json:"conversation_id" gorm:"primaryKey"` // 会话ID，主键
	User1ID        string `json:"user1_id"`                          // 用户1的ID
	User2ID        string `json:"user2_id"`                          // 用户2的ID
	LastMessageID  string `json:"last_message_id"`                   // 最后一条消息ID
}
