package models

import (
	"gorm.io/gorm"
)

// Message 消息模型
// 消息表存储用户之间的消息
type Message struct {
	gorm.Model
	MessageID      string `json:"message_id" gorm:"primaryKey"` // 消息ID，主键
	ConversationID string `json:"conversation_id"`              // 会话ID
	SenderID       string `json:"sender_id"`                    // 发送者用户ID
	ReceiverID     string `json:"receiver_id"`                  // 接收者用户ID
	Content        string `json:"content"`                      // 消息内容
	MessageType    string `json:"message_type"`                 // 消息类型（文本、图片、视频等）
	Status         string `json:"status"`                       // 消息状态（已发送、已读、未读等）
}
