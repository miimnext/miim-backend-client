package services

import (
	"go_core/config"
	"go_core/models"

	"gorm.io/gorm"
)

// 创建评论
func CreateComment(content string, postID uint, authorID uint) (*models.Comment, error) {

	// 创建评论
	comment := models.Comment{
		Content:  content,
		PostID:   postID,
		AuthorID: authorID,
	}

	// 保存评论到数据库
	err := config.DB.Create(&comment).Error
	if err != nil {
		return nil, err
	}

	// 更新帖子评论数量
	var post models.Post
	err = config.DB.Model(&post).Where("id = ?", postID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// 获取帖子所有评论
func GetCommentsByPostID(postID uint) ([]models.Comment, error) {

	var comments []models.Comment
	err := config.DB.Preload("Author").Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// 更新评论内容
func UpdateComment(commentID uint, newContent string) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}

	comment.Content = newContent
	err = config.DB.Save(&comment).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// 删除评论
func DeleteComment(commentID uint, postID uint) error {
	var comment models.Comment
	err := config.DB.First(&comment, commentID).Error
	if err != nil {
		return err
	}

	// 删除评论
	err = config.DB.Delete(&comment).Error
	if err != nil {
		return err
	}

	// 更新帖子评论数量
	var post models.Post
	err = config.DB.Model(&post).Where("id = ?", postID).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}
