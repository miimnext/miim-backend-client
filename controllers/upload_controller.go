package controllers

import (
	"fmt"
	"go_core/utils"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondFailed(c, "Invalid file")
		return
	}

	// 保存文件到服务器
	filePath := "./static/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.RespondFailed(c, "Failed to save file")
		return
	}

	// 获取当前请求的域名和端口
	host := c.Request.Host
	// 构建完整的 URL
	fileURL := fmt.Sprintf("http://%s/static/%s", host, file.Filename)

	utils.RespondSuccess(c, gin.H{
		"filePath": fileURL,
	}, nil)
}
