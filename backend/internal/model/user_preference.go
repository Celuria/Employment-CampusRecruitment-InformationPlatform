package model

import "time"

// UserPreference 用户偏好表
type UserPreference struct {
	ID                 uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             uint64      `gorm:"uniqueIndex:uk_user_preferences_user_id;not null" json:"userId"`
	TargetPositions    JSONStrings `gorm:"type:json" json:"targetPositions"`
	PreferredCities    JSONStrings `gorm:"type:json" json:"preferredCities"`
	PreferredCompanies JSONStrings `gorm:"type:json" json:"preferredCompanies"`
	FocusCompanies     JSONStrings `gorm:"type:json" json:"focusCompanies"`
	RemindBefore       JSONStrings `gorm:"type:json" json:"remindBefore"`
	CreatedAt          time.Time   `json:"createdAt"`
	UpdatedAt          time.Time   `json:"updatedAt"`
}

func (UserPreference) TableName() string { return "user_preferences" }
