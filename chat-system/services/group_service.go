package services

import (
	"chat-system/config"
	"chat-system/models"
	"fmt"
	"time"
)

// 创建群组
func createGroup(ownerID, groupName, description string) (*models.Group, error) {
	group := models.Group{
		GroupID:     fmt.Sprintf("%s-%d", ownerID, time.Now().UnixNano()),
		GroupName:   groupName,
		OwnerID:     ownerID,
		Description: description,
	}

	if err := config.DB.Create(&group).Error; err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return &group, nil
}

// 用户加入群组
func joinGroup(userID, groupID string) error {
	var group models.Group
	if err := config.DB.Where("group_id = ?", groupID).First(&group).Error; err != nil {
		return fmt.Errorf("group not found: %v", err)
	}

	groupMember := models.GroupMember{
		GroupID: groupID,
		UserID:  userID,
	}

	if err := config.DB.Create(&groupMember).Error; err != nil {
		return fmt.Errorf("failed to add user to group: %v", err)
	}

	return nil
}
