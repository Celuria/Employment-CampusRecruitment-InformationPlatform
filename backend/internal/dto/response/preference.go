package response

import "time"

// UserPreferenceResponse 用户偏好响应
type UserPreferenceResponse struct {
	TargetPositions    []string `json:"targetPositions"`
	PreferredCities    []string `json:"preferredCities"`
	PreferredCompanies []string `json:"preferredCompanies"`
	FocusCompanies     []string `json:"focusCompanies"`
	RemindBefore       []string `json:"remindBefore"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
}
