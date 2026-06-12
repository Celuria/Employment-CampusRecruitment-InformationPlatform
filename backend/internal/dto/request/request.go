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

// AdminCareerTalkQuery 管理端宣讲会查询
type AdminCareerTalkQuery struct {
	CareerTalkQuery
	PublishStatus string `form:"publishStatus"`
	SourceType    string `form:"sourceType"`
}

// AdminCareerTalkCreateRequest 创建宣讲会
type AdminCareerTalkCreateRequest struct {
	Title           string   `json:"title" binding:"required,max=200"`
	Company         string   `json:"company" binding:"required,max=128"`
	IndustryCode    string   `json:"industryCode"`
	CompanySize     string   `json:"companySize"`
	StartTime       string   `json:"startTime" binding:"required"`
	EndTime         string   `json:"endTime"`
	Location        string   `json:"location"`
	Campus          string   `json:"campus" binding:"required,oneof=nanhu mafangshan yujiato online"`
	Venue           string   `json:"venue" binding:"required,max=128"`
	Format          string   `json:"format" binding:"required,oneof=online offline hybrid"`
	Positions       []string `json:"positions"`
	TargetMajors    []string `json:"targetMajors"`
	RegistrationURL string   `json:"registrationUrl"`
	SourceURL       string   `json:"sourceUrl"`
	LogoURL         string   `json:"logoUrl"`
	Description     string   `json:"description"`
	PublishStatus   string   `json:"publishStatus"`
}

// AdminCareerTalkUpdateRequest 更新宣讲会
type AdminCareerTalkUpdateRequest struct {
	Title           *string  `json:"title"`
	Company         *string  `json:"company"`
	IndustryCode    *string  `json:"industryCode"`
	CompanySize     *string  `json:"companySize"`
	StartTime       *string  `json:"startTime"`
	EndTime         *string  `json:"endTime"`
	Campus          *string  `json:"campus"`
	Venue           *string  `json:"venue"`
	Format          *string  `json:"format"`
	Positions       []string `json:"positions"`
	TargetMajors    []string `json:"targetMajors"`
	RegistrationURL *string  `json:"registrationUrl"`
	SourceURL       *string  `json:"sourceUrl"`
	LogoURL         *string  `json:"logoUrl"`
	Description     *string  `json:"description"`
	PublishStatus   *string  `json:"publishStatus"`
}

// AdminJobFairQuery 管理端双选会查询
type AdminJobFairQuery struct {
	JobFairQuery
	PublishStatus string `form:"publishStatus"`
	SourceType    string `form:"sourceType"`
}

// AdminJobFairCreateRequest 创建双选会
type AdminJobFairCreateRequest struct {
	Title          string   `json:"title" binding:"required,max=200"`
	StartDate      string   `json:"startDate" binding:"required"`
	EndDate        string   `json:"endDate"`
	StartTime      string   `json:"startTime"`
	Campus         string   `json:"campus" binding:"required,oneof=nanhu mafangshan yujiato online"`
	Venue          string   `json:"venue" binding:"required,max=128"`
	CompanyCount   *int     `json:"companyCount"`
	TargetAudience string   `json:"targetAudience"`
	TargetMajors   []string `json:"targetMajors"`
	Deadline       string   `json:"deadline"`
	DetailURL      string   `json:"detailUrl"`
	SourceURL      string   `json:"sourceUrl"`
	Description    string   `json:"description"`
	PublishStatus  string   `json:"publishStatus"`
}

// AdminJobFairUpdateRequest 更新双选会
type AdminJobFairUpdateRequest struct {
	Title          *string  `json:"title"`
	StartDate      *string  `json:"startDate"`
	EndDate        *string  `json:"endDate"`
	StartTime      *string  `json:"startTime"`
	Campus         *string  `json:"campus"`
	Venue          *string  `json:"venue"`
	CompanyCount   *int     `json:"companyCount"`
	TargetAudience *string  `json:"targetAudience"`
	TargetMajors   []string `json:"targetMajors"`
	Deadline       *string  `json:"deadline"`
	DetailURL      *string  `json:"detailUrl"`
	SourceURL      *string  `json:"sourceUrl"`
	Description    *string  `json:"description"`
	PublishStatus  *string  `json:"publishStatus"`
}

// BatchPublishStatusRequest 批量更新发布状态
type BatchPublishStatusRequest struct {
	IDs           []uint64 `json:"ids" binding:"required,min=1"`
	PublishStatus string   `json:"publishStatus" binding:"required,oneof=draft published archived"`
}

// AdminUserQuery 用户列表查询
type AdminUserQuery struct {
	Keyword  string `form:"keyword"`
	Role     string `form:"role"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// AdminUserCreateRequest 创建用户
type AdminUserCreateRequest struct {
	Username string `json:"username" binding:"required,min=4,max=32"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=student admin"`
}

// AdminUserStatusRequest 更新用户状态
type AdminUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active locked disabled"`
}

// AdminUserUpdateRequest 更新用户资料与角色
type AdminUserUpdateRequest struct {
	Name    *string `json:"name"`
	Email   *string `json:"email"`
	College *string `json:"college"`
	Major   *string `json:"major"`
	Role    *string `json:"role"`
}

// AdminResetPasswordRequest 重置密码
type AdminResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// AuditLogQuery 审计日志查询
type AuditLogQuery struct {
	OperatorID   uint64 `form:"operatorId"`
	Action       string `form:"action"`
	ResourceType string `form:"resourceType"`
	StartDate    string `form:"startDate"`
	EndDate      string `form:"endDate"`
	Page         int    `form:"page"`
	PageSize     int    `form:"pageSize"`
}

// SyncLogQuery 同步记录查询
type SyncLogQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}
