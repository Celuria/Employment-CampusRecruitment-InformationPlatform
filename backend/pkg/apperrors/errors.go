package apperrors

import "fmt"

// AppError 业务错误，携带 HTTP 状态码与业务 code
type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func New(httpStatus, code int, message string) *AppError {
	return &AppError{HTTPStatus: httpStatus, Code: code, Message: message}
}

// 通用错误
var (
	ErrBadRequest      = New(400, 40000, "请求参数错误")
	ErrUnauthorized    = New(401, 40100, "未登录或 Token 已过期")
	ErrForbidden       = New(403, 40300, "无权限访问")
	ErrNotFound        = New(404, 40400, "资源不存在")
	ErrConflict        = New(409, 40900, "资源冲突")
	ErrAccountLocked   = New(423, 42300, "账号已锁定，请稍后再试")
	ErrInternalServer  = New(500, 50000, "服务器内部错误")
)

// 认证模块
var (
	ErrInvalidCredentials = New(401, 40101, "账号或密码错误")
	ErrInvalidToken       = New(401, 40102, "Token 无效")
	ErrTokenExpired       = New(401, 40103, "Token 已过期")
	ErrUsernameExists     = New(400, 40001, "用户名已存在")
	ErrEmailExists        = New(400, 40002, "邮箱已被注册")
	ErrWeakPassword       = New(400, 40003, "密码不符合复杂度要求")
	ErrProfileIncomplete  = New(400, 40201, "资料必填项不完整，请至少填写一个意向岗位")
)

// 日历模块
var (
	ErrCalendarDuplicate = New(409, 40901, "该活动已在您的日历中")
	ErrEventNotFound     = New(404, 40402, "关联活动不存在")
	ErrCalendarNotFound  = New(404, 40403, "日历事件不存在")
)

// 管理端
var (
	ErrAdminRequired     = New(403, 40301, "需要管理员权限")
	ErrCareerTalkNotFound = New(404, 40401, "宣讲会不存在")
	ErrJobFairNotFound   = New(404, 40404, "双选会不存在")
	ErrUserNotFound      = New(404, 40405, "用户不存在")
)
