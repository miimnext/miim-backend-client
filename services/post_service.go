package services

import (
	"fmt"
	"go_core/config"
	"go_core/models"
	"go_core/utils"

	"gorm.io/gorm"
)

type PostService struct{}

// CreatePost 用于创建文章
func (s *PostService) CreatePost(title string, content string, authorID uint, categoryIDs []uint, tagIDs []uint) error {
	// 获取 Author
	var author models.User
	if err := config.DB.First(&author, authorID).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to fetch author with ID %d", authorID), err)
		return fmt.Errorf("failed to fetch author: %w", err)
	}

	// 获取 Categories
	var categories []models.Category
	if err := config.DB.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to fetch categories with IDs %v", categoryIDs), err)
		return fmt.Errorf("failed to fetch categories: %w", err)
	}

	// 获取 Tags
	var tags []models.Tag
	if err := config.DB.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to fetch tags with IDs %v", tagIDs), err)
		return fmt.Errorf("failed to fetch tags: %w", err)
	}

	// 创建 Post
	post := models.Post{
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		Categories: categories,
		Tags:       tags,
	}

	// 保存到数据库
	if err := config.DB.Create(&post).Error; err != nil {
		utils.LogError("Failed to create post", err)
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetPosts 用于获取文章列表
func (s *PostService) GetPosts(pagination *utils.Pagination) ([]models.Post, int64, error) {
	var posts []models.Post
	var totalPosts int64
	// 构建查询
	query := config.DB.Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nickname, username")
	}).Order("created_at DESC")

	// 获取总数
	if err := query.Model(&models.Post{}).Count(&totalPosts).Error; err != nil {
		utils.LogError("Failed to count posts", err)
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// 获取分页数据
	offset, limit := pagination.Paginate()
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		utils.LogError("Failed to fetch posts", err)
		return nil, 0, fmt.Errorf("failed to fetch posts: %w", err)
	}

	return posts, totalPosts, nil
}

// GetPostByID 用于根据 ID 获取单篇文章
func (s *PostService) GetPostByID(id string) (*models.Post, error) {
	// 从数据库中查询单篇文章，使用 Preload 来加载关联的 Categories 和 Tags
	var post models.Post
	err := config.DB.Preload("Categories").Preload("Tags").Preload("Author").Where("id = ?", id).First(&post).Error
	if err != nil {
		// 判断是记录不存在错误还是其他错误
		if err == gorm.ErrRecordNotFound {
			// 如果未找到记录，返回 nil 和自定义错误
			return nil, nil
		}
		// 其他错误，返回数据库查询错误
		return nil, err
	}
	return &post, nil
}
