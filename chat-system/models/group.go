package models

import (
	"gorm.io/gorm"
)

// Group 群组模型
// 群组表存储群组的详细信息
type Group struct {
	gorm.Model
	GroupID     string `json:"group_id" gorm:"primaryKey"` // 群组ID，主键
	GroupName   string `json:"group_name"`                 // 群组名称
	OwnerID     string `json:"owner_id"`                   // 群主ID
	Description string `json:"description"`                // 群组描述
}
