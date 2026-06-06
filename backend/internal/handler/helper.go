package handler

import (
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

// abortError 中断请求并记录业务错误
func abortError(c *gin.Context, err *apperrors.AppError) {
	_ = c.Error(err)
	c.Abort()
}

// bindJSON 绑定 JSON 请求体
func bindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return false
	}
	return true
}

// bindQuery 绑定 Query 参数
func bindQuery(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindQuery(req); err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return false
	}
	return true
}
