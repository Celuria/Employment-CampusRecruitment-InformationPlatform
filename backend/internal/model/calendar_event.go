package model

import "time"

type EventType string

const (
	EventTypeCareerTalk EventType = "career_talk"
	EventTypeJobFair    EventType = "job_fair"
)

// CalendarEvent 日历事件表
type CalendarEvent struct {
	ID           uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64      `gorm:"uniqueIndex:uk_calendar_user_event;index:idx_calendar_user_start;not null" json:"-"`
	EventType    EventType   `gorm:"uniqueIndex:uk_calendar_user_event;size:16;not null" json:"eventType"`
	RefID        uint64      `gorm:"uniqueIndex:uk_calendar_user_event;not null" json:"refId"`
	Title        string      `gorm:"size:200;not null" json:"title"`
	StartTime    time.Time   `gorm:"index:idx_calendar_user_start;not null" json:"startTime"`
	EndTime      *time.Time  `json:"endTime"`
	Location     string      `gorm:"size:256" json:"location"`
	CustomNote   string      `gorm:"size:500" json:"customNote"`
	RemindBefore JSONStrings `gorm:"type:json" json:"remindBefore"`
	Status       string      `gorm:"size:16;default:active" json:"-"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

func (CalendarEvent) TableName() string { return "calendar_events" }
