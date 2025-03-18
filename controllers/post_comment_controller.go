package controllers

import (
	"go_core/services"
	"go_core/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建评论的控制器
func CreateCommentHandler(c *gin.Context) {
	var request struct {
		Content string `json:"content" binding:"required"`
		PostID  uint   `json:"post_id" binding:"required"`
	}

	// 绑定请求的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		return
	}

	userInfo, ok := user.(*services.Claims)
	if !ok {
		return
	}

	// 调用服务层创建评论
	comment, err := services.CreateComment(request.Content, request.PostID, userInfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回创建的评论数据
	utils.RespondSuccess(c, comment, nil)
}

// 获取帖子的评论列表
func GetCommentsByID(c *gin.Context) {
	postID := c.Param("post_id")

	// 转换 post_id 为 uint 类型
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	// 调用服务层获取评论
	comments, err := services.GetCommentsByPostID(uint(postIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回评论列表
	utils.RespondSuccess(c, gin.H{"list": comments}, nil)
}

// 更新评论的内容
func UpdateCommentHandler(c *gin.Context) {
	commentID := c.Param("comment_id")

	// 绑定请求数据
	var request struct {
		Content string `json:"content" binding:"required"`
	}

	// 绑定请求的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 转换 comment_id 为 uint 类型
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id"})
		return
	}

	// 调用服务层更新评论
	comment, err := services.UpdateComment(uint(commentIDInt), request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的评论数据
	utils.RespondSuccess(c, comment, nil)
}

// 删除评论
func DeleteCommentHandler(c *gin.Context) {
	commentID := c.Param("comment_id")
	postID := c.Param("post_id")

	// 转换 comment_id 和 post_id 为 uint 类型
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id"})
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	// 调用服务层删除评论
	err = services.DeleteComment(uint(commentIDInt), uint(postIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功的响应
	utils.RespondSuccess(c, gin.H{"message": "Comment deleted successfully"}, nil)
}
