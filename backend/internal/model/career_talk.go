package model

import "time"

type EventFormat string

const (
	FormatOnline  EventFormat = "online"
	FormatOffline EventFormat = "offline"
	FormatHybrid  EventFormat = "hybrid"
)

type EventStatus string

const (
	EventUpcoming EventStatus = "upcoming"
	EventEnded    EventStatus = "ended"
)

type PublishStatus string

const (
	PublishDraft     PublishStatus = "draft"
	PublishPublished PublishStatus = "published"
	PublishArchived  PublishStatus = "archived"
)

// CareerTalk 宣讲会表
type CareerTalk struct {
	ID              uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string        `gorm:"size:200;not null" json:"title"`
	Company         string        `gorm:"size:128;not null" json:"company"`
	IndustryCode    string        `gorm:"size:32" json:"industryCode"`
	CompanySize     string        `gorm:"size:64" json:"companySize"`
	StartTime       time.Time     `gorm:"not null;index:idx_career_talks_start_time" json:"startTime"`
	EndTime         *time.Time    `json:"endTime"`
	Location        string        `gorm:"size:256;not null" json:"location"`
	Campus          string        `gorm:"size:32" json:"campus"`
	Venue           string        `gorm:"size:128" json:"venue"`
	Format          EventFormat   `gorm:"size:16;not null" json:"format"`
	Positions       JSONStrings   `gorm:"type:json" json:"positions"`
	TargetMajors    JSONStrings   `gorm:"type:json" json:"targetMajors"`
	RegistrationURL string        `gorm:"size:512" json:"registrationUrl"`
	SourceURL       string        `gorm:"size:512" json:"sourceUrl"`
	LogoURL         string        `gorm:"size:512" json:"logoUrl"`
	Description     string        `gorm:"type:text" json:"description"`
	PublishStatus   PublishStatus `gorm:"size:16;default:draft;index:idx_career_talks_publish_status" json:"publishStatus"`
	SourceType      string        `gorm:"size:16;default:manual" json:"sourceType"`
	SyncedAt        *time.Time    `json:"syncedAt"`
	CreatedBy       *uint64       `json:"createdBy"`
	UpdatedBy       *uint64       `json:"updatedBy"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	InCalendar      bool          `json:"inCalendar,omitempty" gorm:"-"`
}

func (CareerTalk) TableName() string { return "career_talks" }
