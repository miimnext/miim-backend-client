package models

import (
	"time"

	"gorm.io/gorm"
)

type ReactionType string

const (
	Like    ReactionType = "like"
	Dislike ReactionType = "dislike"
	None    ReactionType = "none"
)

type PostReaction struct {
	ID        uint         `gorm:"primaryKey"`
	PostID    uint         `gorm:"index:idx_post_user;not null"` // Composite index on PostID and UserID
	UserID    uint         `gorm:"index:idx_post_user;not null"`
	Type      ReactionType `gorm:"type:enum('like','dislike','none');default:'none'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
