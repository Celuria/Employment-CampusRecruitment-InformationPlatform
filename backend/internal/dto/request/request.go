package request

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=4,max=32"`
	Password   string `json:"password" binding:"required,min=8"`
	Email      string `json:"email" binding:"required,email"`
	Captcha    string `json:"captcha"`
	CaptchaKey string `json:"captchaKey"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
	Captcha  string `json:"captcha"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Name            string   `json:"name" binding:"required"`
	College         string   `json:"college" binding:"required"`
	Major           string   `json:"major" binding:"required"`
	Grade           string   `json:"grade"`
	TargetPositions []string `json:"targetPositions" binding:"required,min=1"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email" binding:"required,email"`
}

// UpdatePreferenceRequest 更新偏好请求
type UpdatePreferenceRequest struct {
	TargetPositions    []string `json:"targetPositions"`
	PreferredCities    []string `json:"preferredCities"`
	PreferredCompanies []string `json:"preferredCompanies"`
	FocusCompanies     []string `json:"focusCompanies"`
	RemindBefore       []string `json:"remindBefore"`
}

// CareerTalkQuery 宣讲会查询
type CareerTalkQuery struct {
	Keyword   string `form:"keyword"`
	DateRange string `form:"dateRange"`
	Campus    string `form:"campus"`
	Industry  string `form:"industry"`
	Company   string `form:"company"`
	SortBy    string `form:"sortBy"`
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
}

// JobFairQuery 双选会查询
type JobFairQuery struct {
	Keyword   string `form:"keyword"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
	Campus    string `form:"campus"`
	SortBy    string `form:"sortBy"`
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
}

// CreateCalendarEventRequest 添加日历
type CreateCalendarEventRequest struct {
	EventType    string   `json:"eventType" binding:"required,oneof=career_talk job_fair"`
	RefID        uint64   `json:"refId" binding:"required"`
	CustomNote   string   `json:"customNote"`
	RemindBefore []string `json:"remindBefore"`
}

// UpdateCalendarEventRequest 更新日历
type UpdateCalendarEventRequest struct {
	CustomNote   string   `json:"customNote"`
	RemindBefore []string `json:"remindBefore"`
}

// DismissRecommendationRequest 不感兴趣
type DismissRecommendationRequest struct {
	EventType string `json:"eventType" binding:"required,oneof=career_talk job_fair"`
}

// SyncTriggerRequest 触发同步
type SyncTriggerRequest struct {
	SourceType string `json:"sourceType"`
	Force      bool   `json:"force"`
}
