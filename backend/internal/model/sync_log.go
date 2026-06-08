package model

import "time"

// SyncLog 信息同步任务记录
type SyncLog struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID       string     `gorm:"size:64;not null;index:idx_sync_task" json:"taskId"`
	SourceType   string     `gorm:"size:32;not null" json:"sourceType"`
	Status       string     `gorm:"size:16;not null" json:"status"`
	AddedCount   int        `gorm:"default:0" json:"addedCount"`
	UpdatedCount int        `gorm:"default:0" json:"updatedCount"`
	FailedCount  int        `gorm:"default:0" json:"failedCount"`
	StartedAt    time.Time  `gorm:"not null" json:"startedAt"`
	FinishedAt   *time.Time `json:"finishedAt"`
	OperatorID   uint64     `gorm:"not null" json:"operatorId"`
	ErrorMessage string     `gorm:"size:512" json:"errorMessage"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (SyncLog) TableName() string { return "sync_logs" }
