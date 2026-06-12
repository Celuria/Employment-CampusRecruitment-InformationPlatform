package model

import "time"

// JobFair 双选会表
type JobFair struct {
	ID              uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string        `gorm:"size:200;not null" json:"title"`
	StartDate       time.Time     `gorm:"type:date;not null;index:idx_job_fairs_start_date" json:"startDate"`
	EndDate         *time.Time    `gorm:"type:date" json:"endDate"`
	StartTime       *time.Time    `json:"startTime"`
	Location        string        `gorm:"size:256;not null" json:"location"`
	Campus          string        `gorm:"size:32" json:"campus"`
	Venue           string        `gorm:"size:128" json:"venue"`
	CompanyCount    *int          `json:"companyCount"`
	TargetAudience  string        `gorm:"size:256" json:"targetAudience"`
	TargetMajors    JSONStrings   `gorm:"type:json" json:"targetMajors"`
	Deadline        *time.Time    `gorm:"index:idx_job_fairs_deadline" json:"deadline"`
	DetailURL       string        `gorm:"size:512" json:"detailUrl"`
	SourceURL       string        `gorm:"size:512" json:"sourceUrl"`
	Description     string        `gorm:"type:text" json:"description"`
	PublishStatus   PublishStatus `gorm:"size:16;default:draft;index:idx_job_fairs_publish_status" json:"publishStatus"`
	SourceType      string        `gorm:"size:16;default:manual" json:"sourceType"`
	SyncedAt        *time.Time    `json:"syncedAt"`
	CreatedBy       *uint64       `json:"createdBy"`
	UpdatedBy       *uint64       `json:"updatedBy"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	InCalendar      bool          `json:"inCalendar,omitempty" gorm:"-"`
}

func (JobFair) TableName() string { return "job_fairs" }
