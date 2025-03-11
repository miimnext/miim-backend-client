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
func (s *PostService) CreatePost(title string, content string, authorID uint, CategoryID uint, tagIDs []uint) error {
	// 获取 Author
	var author models.User
	if err := config.DB.First(&author, authorID).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to fetch author with ID %d", authorID), err)
		return fmt.Errorf("failed to fetch author: %w", err)
	}

	// 获取 Category
	var Category models.Category
	if err := config.DB.First(&Category, CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.LogError(fmt.Sprintf("Category with ID %d not found", CategoryID), err)
			return fmt.Errorf("category with ID %d not found", CategoryID)
		}
		utils.LogError(fmt.Sprintf("Failed to fetch category with ID %d", CategoryID), err)
		return fmt.Errorf("failed to fetch category: %w", err)
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
		CategoryID: CategoryID,
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
	}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nickname, username,avatar")
	}).Order("created_at DESC, id DESC")

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
	// 从数据库中查询单篇文章，使用 Preload 来加载关联的 Category 和 Tags
	var post models.Post
	err := config.DB.Preload("Category").Preload("Tags").Preload("Author").Where("id = ?", id).First(&post).Error
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

func (s *PostService) GetPostsByUser(pagination *utils.Pagination, id string) ([]models.Post, int64, error) {
	var posts []models.Post
	var totalPosts int64
	// 打印 SQL 语句
	// 查询该用户名的所有帖子
	query := config.DB.Preload("Category").Preload("Tags").Preload("Author").Where("author_id= ?", id).Order("created_at DESC, id DESC")
	fmt.Println(query.Statement.SQL.String()) // 打印 SQL 语句

	// 获取总数
	if err := query.Model(&models.Post{}).Count(&totalPosts).Error; err != nil {
		utils.LogError("Failed to count posts", err)
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// 获取分页数据
	offset, limit := pagination.Paginate()
	// 设置最大 limit，防止查询过多数据
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		utils.LogError("Failed to fetch posts", err)
		return nil, 0, fmt.Errorf("failed to fetch posts: %w", err)
	}

	return posts, totalPosts, nil
}

// DeletePost 根据 ID 删除文章
func (s *PostService) DeletePost(id uint) error {
	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("post with ID %d not found", id)
		}
		utils.LogError(fmt.Sprintf("Failed to find post with ID %d", id), err)
		return fmt.Errorf("failed to find post: %w", err)
	}

	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		utils.LogError(fmt.Sprintf("Failed to delete post with ID %d", id), err)
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
