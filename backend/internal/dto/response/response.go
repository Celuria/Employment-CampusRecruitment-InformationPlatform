package response

import "time"

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
	TokenType string `json:"tokenType"`
}

// UserProfileResponse 用户资料
type UserProfileResponse struct {
	ID               uint64    `json:"id"`
	Username         string    `json:"username"`
	Role             string    `json:"role"`
	Status           string    `json:"status"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	College          string    `json:"college"`
	Major            string    `json:"major"`
	Grade            string    `json:"grade"`
	TargetPositions  []string  `json:"targetPositions"`
	Phone            string    `json:"phone"`
	Avatar           string    `json:"avatar"`
	ProfileCompleted bool      `json:"profileCompleted"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// HealthResponse 健康检查
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
