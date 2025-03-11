package services

import (
	"fmt"
	"go_core/config"
	"go_core/models"

	"gorm.io/gorm"
)

func PostReaction(postID, userID uint, reactionType models.ReactionType) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查文章是否存在
		var post models.Post
		if err := tx.First(&post, postID).Error; err != nil {
			return fmt.Errorf("failed to find post: %w", err)
		}
		// 检查用户是否已有反应
		var reaction models.PostReaction
		err := tx.Where("post_id = ? AND user_id = ?", postID, userID).First(&reaction).Error
		// 计算点赞数的变化
		likesChange := 0
		switch {
		case err == gorm.ErrRecordNotFound:
			// 如果没有反应记录，创建新记录
			if reactionType == models.Like {
				reaction = models.PostReaction{PostID: postID, UserID: userID, Type: models.Like}
				likesChange = 1
			} else if reactionType == models.Dislike {
				reaction = models.PostReaction{PostID: postID, UserID: userID, Type: models.Dislike}
				likesChange = -1
			}
			if err := tx.Create(&reaction).Error; err != nil {
				return fmt.Errorf("failed to create reaction: %w", err)
			}
		case err != nil:
			return fmt.Errorf("failed to check reaction: %w", err)

		case reaction.Type == models.Like:
			// 已经点赞
			if reactionType == models.Like {
				// 取消点赞 (-1)
				if err := tx.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.PostReaction{}).Error; err != nil {
					return fmt.Errorf("failed to delete reaction: %w", err)
				}
				likesChange = -1
			} else if reactionType == models.Dislike {
				// 点踩 (-2)
				if err := tx.Model(&reaction).Update("type", models.Dislike).Error; err != nil {
					return fmt.Errorf("failed to update reaction: %w", err)
				}
				likesChange = -2
			}

		case reaction.Type == models.Dislike:
			// 已经点踩
			if reactionType == models.Dislike {
				// 取消点踩 (+1)
				if err := tx.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.PostReaction{}).Error; err != nil {
					return fmt.Errorf("failed to delete reaction: %w", err)
				}
				likesChange = 1
			} else if reactionType == models.Like {
				// 点赞 (+2)
				if err := tx.Model(&reaction).Update("type", models.Like).Error; err != nil {
					return fmt.Errorf("failed to update reaction: %w", err)
				}
				likesChange = 2
			}
		}
		if err := tx.Model(&models.Post{}).Where("id = ?", postID).
			Update("likes", gorm.Expr("likes + ?", likesChange)).Error; err != nil {
			return fmt.Errorf("failed to update likes: %w", err)
		}

		return nil
	})
}

// GetPostReactionByUser 检查用户对指定帖子的反应（点赞/点踩）
func GetPostReactionByUser(postID, userID uint) (models.ReactionType, error) {
	var reaction models.PostReaction
	// 查询用户是否对该帖子有反应
	if err := config.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&reaction).Error; err != nil {
		// 如果没有找到记录，返回 "none" 和 nil 错误
		if err == gorm.ErrRecordNotFound {
			return models.None, nil
		}
		// 如果其他错误，返回 err
		return models.None, err
	}
	// 如果找到了记录，返回 reaction 类型（点赞或点踩）
	return reaction.Type, nil
}
