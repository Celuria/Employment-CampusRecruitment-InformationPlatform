package response

// CalendarEventVO 日历事件响应
type CalendarEventVO struct {
	ID             uint64   `json:"id"`
	EventType      string   `json:"eventType"`
	RefID          uint64   `json:"refId"`
	Title          string   `json:"title"`
	StartTime      string   `json:"startTime"`
	EndTime        string   `json:"endTime,omitempty"`
	Location       string   `json:"location"`
	CustomNote     string   `json:"customNote,omitempty"`
	RemindBefore   []string `json:"remindBefore"`
	ReminderStatus string   `json:"reminderStatus,omitempty"`
	SourceUpdated  bool     `json:"sourceUpdated,omitempty"`
	CreatedAt      string   `json:"createdAt"`
}
