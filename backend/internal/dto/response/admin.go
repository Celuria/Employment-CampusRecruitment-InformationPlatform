package response

import (
	"time"

	"github.com/employment-center/campus-recruitment/internal/model"
)

// AdminCareerTalkVO 管理端宣讲会
type AdminCareerTalkVO struct {
	ID              uint64  `json:"id"`
	Title           string  `json:"title"`
	Company         string  `json:"company"`
	IndustryCode    string  `json:"industryCode"`
	CompanySize     string  `json:"companySize"`
	StartTime       string  `json:"startTime"`
	EndTime         *string `json:"endTime,omitempty"`
	Location        string  `json:"location"`
	Campus          string  `json:"campus"`
	Venue           string  `json:"venue"`
	Format          string  `json:"format"`
	Positions       []string `json:"positions"`
	TargetMajors    []string `json:"targetMajors"`
	RegistrationURL string  `json:"registrationUrl"`
	SourceURL       string  `json:"sourceUrl"`
	LogoURL         string  `json:"logoUrl"`
	Description     string  `json:"description"`
	PublishStatus   string  `json:"publishStatus"`
	SourceType      string  `json:"sourceType"`
	CreatedBy       *uint64 `json:"createdBy"`
	UpdatedBy       *uint64 `json:"updatedBy"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// AdminJobFairVO 管理端双选会
type AdminJobFairVO struct {
	ID             uint64   `json:"id"`
	Title          string   `json:"title"`
	StartDate      string   `json:"startDate"`
	EndDate        *string  `json:"endDate,omitempty"`
	StartTime      *string  `json:"startTime,omitempty"`
	Location       string   `json:"location"`
	Campus         string   `json:"campus"`
	Venue          string   `json:"venue"`
	CompanyCount   *int     `json:"companyCount"`
	TargetAudience string   `json:"targetAudience"`
	TargetMajors   []string `json:"targetMajors"`
	Deadline       *string  `json:"deadline,omitempty"`
	DetailURL      string   `json:"detailUrl"`
	SourceURL      string   `json:"sourceUrl"`
	Description    string   `json:"description"`
	PublishStatus  string   `json:"publishStatus"`
	SourceType     string   `json:"sourceType"`
	CreatedBy      *uint64  `json:"createdBy"`
	UpdatedBy      *uint64  `json:"updatedBy"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

// AdminUserVO 管理端用户
type AdminUserVO struct {
	ID          uint64  `json:"id"`
	Username    string  `json:"username"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	College     string  `json:"college"`
	Major       string  `json:"major"`
	Role        string  `json:"role"`
	Status      string  `json:"status"`
	LastLoginAt *string `json:"lastLoginAt,omitempty"`
	CreatedAt   string  `json:"createdAt"`
}

// SyncLogVO 同步记录
type SyncLogVO struct {
	ID           uint64  `json:"id"`
	TaskID       string  `json:"taskId"`
	SourceType   string  `json:"sourceType"`
	Status       string  `json:"status"`
	AddedCount   int     `json:"addedCount"`
	UpdatedCount int     `json:"updatedCount"`
	FailedCount  int     `json:"failedCount"`
	StartedAt    string  `json:"startedAt"`
	FinishedAt   *string `json:"finishedAt,omitempty"`
	OperatorID   uint64  `json:"operatorId"`
	ErrorMessage string  `json:"errorMessage"`
}

// AuditLogVO 审计日志
type AuditLogVO struct {
	ID           uint64 `json:"id"`
	OperatorID   uint64 `json:"operatorId"`
	OperatorName string `json:"operatorName"`
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
	ResourceID   uint64 `json:"resourceId"`
	Detail       string `json:"detail"`
	IP           string `json:"ip"`
	CreatedAt    string `json:"createdAt"`
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func formatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func ToAdminCareerTalkVO(t *model.CareerTalk) AdminCareerTalkVO {
	return AdminCareerTalkVO{
		ID:              t.ID,
		Title:           t.Title,
		Company:         t.Company,
		IndustryCode:    t.IndustryCode,
		CompanySize:     t.CompanySize,
		StartTime:       t.StartTime.Format(time.RFC3339),
		EndTime:         formatTimePtr(t.EndTime),
		Location:        t.Location,
		Campus:          t.Campus,
		Venue:           t.Venue,
		Format:          string(t.Format),
		Positions:       t.Positions,
		TargetMajors:    t.TargetMajors,
		RegistrationURL: t.RegistrationURL,
		SourceURL:       t.SourceURL,
		LogoURL:         t.LogoURL,
		Description:     t.Description,
		PublishStatus:   string(t.PublishStatus),
		SourceType:      t.SourceType,
		CreatedBy:       t.CreatedBy,
		UpdatedBy:       t.UpdatedBy,
		CreatedAt:       t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       t.UpdatedAt.Format(time.RFC3339),
	}
}

func ToAdminJobFairVO(f *model.JobFair) AdminJobFairVO {
	return AdminJobFairVO{
		ID:             f.ID,
		Title:          f.Title,
		StartDate:      f.StartDate.Format("2006-01-02"),
		EndDate:        formatDatePtr(f.EndDate),
		StartTime:      formatTimePtr(f.StartTime),
		Location:       f.Location,
		Campus:         f.Campus,
		Venue:          f.Venue,
		CompanyCount:   f.CompanyCount,
		TargetAudience: f.TargetAudience,
		TargetMajors:   f.TargetMajors,
		Deadline:       formatTimePtr(f.Deadline),
		DetailURL:      f.DetailURL,
		SourceURL:      f.SourceURL,
		Description:    f.Description,
		PublishStatus:  string(f.PublishStatus),
		SourceType:     f.SourceType,
		CreatedBy:      f.CreatedBy,
		UpdatedBy:      f.UpdatedBy,
		CreatedAt:      f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      f.UpdatedAt.Format(time.RFC3339),
	}
}

func ToAdminUserVO(u *model.User) AdminUserVO {
	return AdminUserVO{
		ID:          u.ID,
		Username:    u.Username,
		Name:        u.Name,
		Email:       u.Email,
		College:     u.College,
		Major:       u.Major,
		Role:        string(u.Role),
		Status:      string(u.Status),
		LastLoginAt: formatTimePtr(u.LastLoginAt),
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
	}
}

func ToSyncLogVO(l *model.SyncLog) SyncLogVO {
	return SyncLogVO{
		ID:           l.ID,
		TaskID:       l.TaskID,
		SourceType:   l.SourceType,
		Status:       l.Status,
		AddedCount:   l.AddedCount,
		UpdatedCount: l.UpdatedCount,
		FailedCount:  l.FailedCount,
		StartedAt:    l.StartedAt.Format(time.RFC3339),
		FinishedAt:   formatTimePtr(l.FinishedAt),
		OperatorID:   l.OperatorID,
		ErrorMessage: l.ErrorMessage,
	}
}

func ToAuditLogVO(l *model.AuditLog, operatorName string) AuditLogVO {
	return AuditLogVO{
		ID:           l.ID,
		OperatorID:   l.OperatorID,
		OperatorName: operatorName,
		Action:       l.Action,
		ResourceType: l.ResourceType,
		ResourceID:   l.ResourceID,
		Detail:       l.Detail,
		IP:           l.IP,
		CreatedAt:    l.CreatedAt.Format(time.RFC3339),
	}
}
