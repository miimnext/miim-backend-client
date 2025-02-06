package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string `json:"content"`
	PostID   uint   `json:"post_id"`
	Post     Post   `json:"post"`
	AuthorID uint   `json:"author_id"`
	Author   User   `json:"author"`
}
