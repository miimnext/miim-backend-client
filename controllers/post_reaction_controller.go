package controllers

import (
	"go_core/models"
	"go_core/services"
	"go_core/utils"

	"github.com/gin-gonic/gin"
)

// 请求体结构体
type PostReactionRequest struct {
	PostID       uint `json:"post_id"`
	UserID       uint `json:"user_id"`
	ReactionType int  `json:"reaction_type"` // 1 表示点赞，-1 表示取消点赞
}

// 用户点赞和取消点赞
func PostReaction(c *gin.Context) {
	var request PostReactionRequest
	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.RespondFailed(c, "Invalid request body")
		return
	}
	reactionType := models.None
	if request.ReactionType == 1 {
		reactionType = models.Like
	} else if request.ReactionType == -1 {
		reactionType = models.Dislike
	}
	// 调用服务层处理点赞/取消点赞
	err := services.PostReaction(request.PostID, request.UserID, reactionType)
	if err != nil {
		utils.RespondFailed(c, err.Error())
		return
	}

	// 返回成功信息
	if request.ReactionType == 1 {
		utils.RespondSuccess(c, gin.H{"message": "Post liked successfully"}, nil)
	} else {
		utils.RespondSuccess(c, gin.H{"message": "Post unliked successfully"}, nil)
	}
}
