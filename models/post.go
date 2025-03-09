package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	AuthorID   uint           `json:"-"`
	Author     User           `gorm:"foreignKey:AuthorID" json:"author"`
	Categories []Category     `gorm:"many2many:post_categories;" json:"categories"`
	Tags       []Tag          `gorm:"many2many:post_tags;" json:"tags"`
	Image      string         `json:"image" gorm:"default:''"` // 添加 Image 字段，默认值为空字符串
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"` // 不返回 DeletedAt
}

type Tag struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Posts     []Post         `gorm:"many2many:post_tags;" json:"-"` // 不返回 Posts
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 不返回 DeletedAt
}

type Category struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Posts     []Post         `gorm:"many2many:post_categories;" json:"-"` // 不返回 Posts
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 不返回 DeletedAt
}
