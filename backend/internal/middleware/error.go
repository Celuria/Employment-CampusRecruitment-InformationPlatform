package middleware

import (
	"errors"

	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/logger"
	"github.com/employment-center/campus-recruitment/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Fail(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
			return
		}

		logger.Log.Error("internal error", zap.Error(err))
		response.Fail(c, 500, apperrors.ErrInternalServer.Code, apperrors.ErrInternalServer.Message)
	}
}

// AuthErrorHandler 处理鉴权/权限 abort 标记
func AuthErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.IsAborted() && c.Writer.Status() == 200 {
			if _, ok := c.Get("forbidden"); ok {
				response.Fail(c, 403, apperrors.ErrForbidden.Code, apperrors.ErrForbidden.Message)
				return
			}
			if _, ok := c.Get("auth_error"); ok {
				response.Fail(c, 401, apperrors.ErrUnauthorized.Code, apperrors.ErrUnauthorized.Message)
			}
		}
	}
}

// Recovery  panic 恢复
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Log.Error("panic recovered", zap.Any("error", recovered))
		response.Fail(c, 500, apperrors.ErrInternalServer.Code, apperrors.ErrInternalServer.Message)
	})
}
