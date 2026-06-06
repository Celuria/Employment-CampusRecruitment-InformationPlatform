package model

import "time"

type ReminderStatus string

const (
	ReminderPending   ReminderStatus = "pending"
	ReminderSent      ReminderStatus = "sent"
	ReminderFailed    ReminderStatus = "failed"
	ReminderCancelled ReminderStatus = "cancelled"
)

// ReminderLog 提醒记录表
type ReminderLog struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CalendarEventID uint64         `gorm:"index:idx_reminder_calendar_event_id;not null" json:"calendarEventId"`
	UserID          uint64         `gorm:"index:idx_reminder_user_id;not null" json:"-"`
	EventTitle      string         `gorm:"size:200" json:"eventTitle"`
	EventType       EventType      `gorm:"size:16" json:"eventType"`
	RemindBefore    string         `gorm:"size:8" json:"remindBefore"`
	ScheduledTime   time.Time      `gorm:"index:idx_reminder_status_scheduled;not null" json:"scheduledTime"`
	SentTime        *time.Time     `json:"sentTime"`
	Status          ReminderStatus `gorm:"size:16;default:pending;index:idx_reminder_status_scheduled" json:"status"`
	RetryCount      int            `gorm:"default:0" json:"retryCount"`
	FailReason      string         `gorm:"size:512" json:"failReason"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

func (ReminderLog) TableName() string { return "reminder_logs" }
