package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	AuthorID     uint           `json:"-"`                                                 // Hidden from JSON response
	ReactionType ReactionType   `gorm:"type:enum('like','dislike','none');default:'none'"` // Enum default value
	Likes        int            `json:"likes" gorm:"default:0"`
	Author       User           `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID   uint           `json:"category_id"`
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category"` // Added Category relation
	Tags         []Tag          `gorm:"many2many:post_tags;" json:"tags"`      // Many-to-many relation with Tags
	Image        string         `json:"image" gorm:"default:''"`               // Default empty string for image
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Hidden DeletedAt from JSON response
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
	Posts     []Post         `gorm:"foreignKey:CategoryID" json:"posts"` // 不返回 Posts
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 不返回 DeletedAt
}
