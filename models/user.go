package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	NickName  string         `gorm:"column:nickname" json:"nickname"`
	Username  string         `json:"username"`
	Password  string         `json:"-"`
	Avatar    string         `gorm:"default:'/static/avatar.jpg'" json:"avatar"`
	Balance   uint           `json:"-"` // Default balance is 0
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
