package models

import (
	"gorm.io/gorm"
)

// User 用户模型
// 用户表存储所有注册用户的信息
type User struct {
	gorm.Model        // gorm.Model 包含 ID、CreatedAt、UpdatedAt 和 DeletedAt 字段
	UserID     string `json:"user_id" gorm:"primaryKey"` // 用户ID，主键
	Username   string `json:"username"`                  // 用户名
	Password   string `json:"password"`                  // 密码
	Email      string `json:"email"`                     // 电子邮箱
	Phone      string `json:"phone"`                     // 手机号码
	AvatarURL  string `json:"avatar_url"`                // 用户头像 URL
	Status     string `json:"status"`                    // 用户状态（在线/离线）
	LastLogin  string `json:"last_login"`                // 最后登录时间
	Bio        string `json:"bio"`                       // 用户简介
}
