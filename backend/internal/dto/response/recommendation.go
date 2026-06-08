package response

// RecommendationVO 个性化推荐项
type RecommendationVO struct {
	ID             uint64   `json:"id"`
	EventType      string   `json:"eventType"`
	RefID          uint64   `json:"refId"`
	Title          string   `json:"title"`
	Company        string   `json:"company,omitempty"`
	StartTime      string   `json:"startTime"`
	StartDate      string   `json:"startDate,omitempty"`
	EndDate        string   `json:"endDate,omitempty"`
	Location       string   `json:"location"`
	Industry       string   `json:"industry,omitempty"`
	IndustryCode   string   `json:"industryCode,omitempty"`
	CompanySize    string   `json:"companySize,omitempty"`
	Format         string   `json:"format,omitempty"`
	Positions      []string `json:"positions,omitempty"`
	Status         string   `json:"status,omitempty"`
	CompanyCount   *int     `json:"companyCount,omitempty"`
	TargetAudience string   `json:"targetAudience,omitempty"`
	Deadline       string   `json:"deadline,omitempty"`
	DetailURL      string   `json:"detailUrl,omitempty"`
	MatchScore     int      `json:"matchScore,omitempty"`
	MatchReasons   []string `json:"matchReasons"`
	InCalendar     bool     `json:"inCalendar"`
}

// RecommendationListResult 推荐列表响应
type RecommendationListResult struct {
	List     []RecommendationVO `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
	Fallback bool               `json:"fallback"`
}
