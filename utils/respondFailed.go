package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RespondFailed(c *gin.Context, message string, code ...int) {
	// 如果传入了 code 参数，就使用它，否则默认 404
	respCode := 404
	if len(code) > 0 && code[0] != 0 {
		respCode = code[0]
	}

	// 构建响应数据
	response := ResponseData{
		Code:    respCode,
		Message: message,
	}

	// 返回响应
	c.JSON(http.StatusOK, response)
}
