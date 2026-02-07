package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构体
type Response struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: ErrCodeSuccess,
		Msg:  ErrorMessages[ErrCodeSuccess],
		Data: data,
	})
}

// ResponseSuccessWithMsg 带自定义消息的成功响应
func ResponseSuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: ErrCodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// ResponseError 错误响应
func ResponseError(c *gin.Context, httpStatus int, code ErrorCode) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  ErrorMessages[code],
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, httpStatus int, code ErrorCode, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
