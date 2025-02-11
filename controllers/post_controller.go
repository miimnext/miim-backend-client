package controllers

import (
	"fmt"
	"go_core/services"
	"go_core/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var postService services.PostService

// CreatePost 用于创建文章
func CreatePost(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		AuthorID    uint   `json:"author_id"`
		CategoryIDs []uint `json:"category_ids"` // 支持多个 CategoryID
		TagIDs      []uint `json:"tag_ids"`      // 支持多个 TagID
	}

	// 绑定输入数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层进行创建
	if err := postService.CreatePost(input.Title, input.Content, input.AuthorID, input.CategoryIDs, input.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功信息
	utils.RespondSuccess(c, nil, nil)
}

// GetAllPosts 用于分页查询所有文章
func GetAllPosts(c *gin.Context) {
	// 获取分页参数
	pagination := utils.GetPagination(c)

	// 查询文章数据
	posts, totalPosts, err := postService.GetPosts(&pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query posts"})
		return
	}

	utils.RespondSuccess(c, gin.H{"list": posts}, &utils.Pagination{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		Total:    totalPosts,
	})

}

// GetPostByID 用于根据 ID 查询单篇文章
func GetPostByID(c *gin.Context) {
	// 获取查询参数中的 id
	id := c.Param("id") // 从路径参数中获取 id

	// 调用服务层根据 id 获取单篇文章
	post, err := postService.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query post"})
		return
	}

	// 如果没有找到文章，则返回空数据
	if post == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Post not found",
			"data":    nil, // 返回空
		})
		return
	}

	// 返回单篇文章
	utils.RespondSuccess(c, post, nil)
}

// DeletePost 处理删除文章的请求
func DeletePost(c *gin.Context) {
	// 从 URL 参数中获取文章 ID
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 调用服务层删除文章
	err = postService.DeletePost(uint(id))
	if err != nil {
		// 文章未找到的情况
		if err.Error() == fmt.Sprintf("post with ID %d not found", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		// 其他错误
		utils.LogError("Failed to delete post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// 删除成功
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
