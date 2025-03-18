package models

import "gorm.io/gorm"

// Comment 模型
type Comment struct {
	gorm.Model
	Content  string `json:"content"`
	PostID   uint   `json:"post_id"`
	AuthorID uint   `json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;" json:"author"`
}
