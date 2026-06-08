package model

import "time"

// AuditLog 管理端操作审计日志
type AuditLog struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OperatorID   uint64    `gorm:"not null;index:idx_audit_operator" json:"operatorId"`
	Action       string    `gorm:"size:16;not null;index:idx_audit_action" json:"action"`
	ResourceType string    `gorm:"size:32;not null;index:idx_audit_resource" json:"resourceType"`
	ResourceID   uint64    `gorm:"default:0" json:"resourceId"`
	Detail       string    `gorm:"type:text" json:"detail"`
	IP           string    `gorm:"size:64" json:"ip"`
	CreatedAt    time.Time `gorm:"index:idx_audit_created" json:"createdAt"`
}

func (AuditLog) TableName() string { return "audit_logs" }
