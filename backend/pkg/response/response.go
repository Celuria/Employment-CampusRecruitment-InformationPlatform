package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 统一 API 响应结构，与前端 ApiResponse 对齐
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageResult 分页响应
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data})
}

// OKMessage 成功响应（自定义 message）
func OKMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: message, Data: data})
}

// Fail 失败响应
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{Code: code, Message: message, Data: nil})
}

// Page 分页成功响应
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	OK(c, PageResult{List: list, Total: total, Page: page, PageSize: pageSize})
}
