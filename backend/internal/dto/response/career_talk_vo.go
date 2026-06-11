package response

import (
	"time"

	"github.com/employment-center/campus-recruitment/internal/model"
)

// CareerTalkVO 宣讲会列表/详情响应（学生端）
type CareerTalkVO struct {
	ID              uint64   `json:"id"`
	Title           string   `json:"title"`
	Company         string   `json:"company"`
	Industry        string   `json:"industry"`
	IndustryCode    string   `json:"industryCode"`
	CompanySize     string   `json:"companySize"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime,omitempty"`
	Location        string   `json:"location"`
	Campus          string   `json:"campus"`
	Format          string   `json:"format"`
	Positions       []string `json:"positions"`
	TargetMajors    []string `json:"targetMajors"`
	RegistrationURL string   `json:"registrationUrl"`
	SourceURL       string   `json:"sourceUrl"`
	LogoURL         string   `json:"logoUrl"`
	Description     string   `json:"description,omitempty"`
	PublishedAt     string   `json:"publishedAt"`
	Status          string   `json:"status"`
	InCalendar      bool     `json:"inCalendar"`
}

// JobFairVO 双选会列表/详情响应（学生端）
type JobFairVO struct {
	ID             uint64   `json:"id"`
	Title          string   `json:"title"`
	StartDate      string   `json:"startDate"`
	EndDate        string   `json:"endDate,omitempty"`
	StartTime      string   `json:"startTime,omitempty"`
	Location       string   `json:"location"`
	Campus         string   `json:"campus"`
	CompanyCount   *int     `json:"companyCount"`
	TargetAudience string   `json:"targetAudience"`
	TargetMajors   []string `json:"targetMajors"`
	Deadline       string   `json:"deadline,omitempty"`
	DetailURL      string   `json:"detailUrl"`
	SourceURL      string   `json:"sourceUrl"`
	Description    string   `json:"description,omitempty"`
	PublishedAt    string   `json:"publishedAt"`
	Status         string   `json:"status"`
	InCalendar     bool     `json:"inCalendar"`
}

var industryLabels = map[string]string{
	"internet":      "互联网",
	"finance":       "金融",
	"manufacturing": "制造",
	"consulting":    "咨询",
}

func industryLabel(code string) string {
	if label, ok := industryLabels[code]; ok {
		return label
	}
	return code
}

func stringsOrEmptySlice(s model.JSONStrings) []string {
	if s == nil {
		return []string{}
	}
	return []string(s)
}

// ToCareerTalkVO 将 model.CareerTalk 转换为学生端 VO
func ToCareerTalkVO(t *model.CareerTalk, inCalendar bool) CareerTalkVO {
	now := time.Now()
	status := "upcoming"
	if t.StartTime.Before(now) {
		status = "ended"
	}
	vo := CareerTalkVO{
		ID:              t.ID,
		Title:           t.Title,
		Company:         t.Company,
		Industry:        industryLabel(t.IndustryCode),
		IndustryCode:    t.IndustryCode,
		CompanySize:     t.CompanySize,
		StartTime:       t.StartTime.Format(time.RFC3339),
		Location:        t.Location,
		Campus:          t.Campus,
		Format:          string(t.Format),
		Positions:       stringsOrEmptySlice(t.Positions),
		TargetMajors:    stringsOrEmptySlice(t.TargetMajors),
		RegistrationURL: t.RegistrationURL,
		SourceURL:       t.SourceURL,
		LogoURL:         t.LogoURL,
		Description:     t.Description,
		PublishedAt:     t.CreatedAt.Format(time.RFC3339),
		Status:          status,
		InCalendar:      inCalendar,
	}
	if t.EndTime != nil {
		vo.EndTime = t.EndTime.Format(time.RFC3339)
	}
	return vo
}

// ToJobFairVO 将 model.JobFair 转换为学生端 VO
func ToJobFairVO(f *model.JobFair, inCalendar bool) JobFairVO {
	vo := JobFairVO{
		ID:             f.ID,
		Title:          f.Title,
		StartDate:      f.StartDate.Format("2006-01-02"),
		Location:       f.Location,
		Campus:         f.Campus,
		CompanyCount:   f.CompanyCount,
		TargetAudience: f.TargetAudience,
		TargetMajors:   stringsOrEmptySlice(f.TargetMajors),
		DetailURL:      f.DetailURL,
		SourceURL:      f.SourceURL,
		Description:    f.Description,
		PublishedAt:    f.CreatedAt.Format(time.RFC3339),
		Status:         "upcoming",
		InCalendar:     inCalendar,
	}
	if f.EndDate != nil {
		vo.EndDate = f.EndDate.Format("2006-01-02")
	}
	if f.StartTime != nil {
		vo.StartTime = f.StartTime.Format(time.RFC3339)
	}
	if f.Deadline != nil {
		vo.Deadline = f.Deadline.Format(time.RFC3339)
	}
	return vo
}