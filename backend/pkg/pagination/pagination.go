package pagination

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Query 分页查询参数
type Query struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

// Normalize 规范化分页参数
func (q *Query) Normalize() (page, pageSize int) {
	page = q.Page
	pageSize = q.PageSize
	if page <= 0 {
		page = DefaultPage
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}

// Offset 计算偏移量
func (q *Query) Offset() int {
	page, pageSize := q.Normalize()
	return (page - 1) * pageSize
}

// Limit 返回 limit
func (q *Query) Limit() int {
	_, pageSize := q.Normalize()
	return pageSize
}

// ParsePage 从字符串解析页码
func ParsePage(s string) int {
	return parseOrDefault(s, DefaultPage)
}

// ParsePageSize 从字符串解析每页条数
func ParsePageSize(s string) int {
	ps := parseOrDefault(s, DefaultPageSize)
	if ps > MaxPageSize {
		return MaxPageSize
	}
	return ps
}

func parseOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return defaultVal
		}
		n = n*10 + int(ch-'0')
	}
	if n <= 0 {
		return defaultVal
	}
	return n
}
