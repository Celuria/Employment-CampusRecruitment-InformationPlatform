package model

import "time"

// UserRole 用户角色
type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleAdmin   UserRole = "admin"
)

// UserStatus 用户状态
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusLocked   UserStatus = "locked"
	StatusDisabled UserStatus = "disabled"
)

// User 用户表
type User struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username        string     `gorm:"size:64;uniqueIndex:uk_users_username;not null" json:"username"`
	PasswordHash    string     `gorm:"size:255;not null" json:"-"`
	Email           string     `gorm:"size:128;uniqueIndex:uk_users_email;not null" json:"email"`
	Name            string     `gorm:"size:64" json:"name"`
	College         string     `gorm:"size:64" json:"college"`
	Major           string     `gorm:"size:64" json:"major"`
	Grade           string     `gorm:"size:16" json:"grade"`
	TargetPositions JSONStrings `gorm:"type:json" json:"targetPositions"`
	Phone           string     `gorm:"size:20" json:"phone"`
	Avatar          string     `gorm:"size:512" json:"avatar"`
	Role            UserRole   `gorm:"size:16;default:student;not null;index:idx_users_role" json:"role"`
	Status          UserStatus `gorm:"size:16;default:active;not null;index:idx_users_status" json:"status"`
	LoginAttempts   int        `gorm:"default:0" json:"-"`
	LockedUntil     *time.Time `json:"-"`
	LastLoginAt     *time.Time `json:"lastLoginAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (User) TableName() string { return "users" }
