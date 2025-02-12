package services

import (
	"chat-system/config"
	"chat-system/models"
	"fmt"
	"time"
)

// 处理消息存储和会话管理
func handleMessage(senderID string, messageType int, message []byte) error {
	messageContent := string(message)
	receiverID := "" // TODO: 获取接收者ID

	msg := models.Message{
		MessageID:      fmt.Sprintf("%s-%d", senderID, time.Now().UnixNano()),
		ConversationID: "",
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Content:        messageContent,
		MessageType:    "text",
		Status:         "sent",
	}

	if err := config.DB.Create(&msg).Error; err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}

	// 处理会话
	if receiverID != "" {
		var conversation models.Conversation
		err := config.DB.Where("user1_id = ? AND user2_id = ?", senderID, receiverID).
			Or("user1_id = ? AND user2_id = ?", receiverID, senderID).
			First(&conversation).Error

		if err != nil {
			conversation = models.Conversation{
				ConversationID: fmt.Sprintf("%s-%s", senderID, receiverID),
				User1ID:        senderID,
				User2ID:        receiverID,
				LastMessageID:  msg.MessageID,
			}
			if err := config.DB.Create(&conversation).Error; err != nil {
				return fmt.Errorf("failed to create conversation: %v", err)
			}
		} else {
			conversation.LastMessageID = msg.MessageID
			if err := config.DB.Save(&conversation).Error; err != nil {
				return fmt.Errorf("failed to update conversation: %v", err)
			}
		}
	}

	return nil
}
