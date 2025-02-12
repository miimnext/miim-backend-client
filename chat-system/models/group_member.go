package models

import (
	"gorm.io/gorm"
)

// GroupMember 群组成员模型
type GroupMember struct {
	gorm.Model
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}
