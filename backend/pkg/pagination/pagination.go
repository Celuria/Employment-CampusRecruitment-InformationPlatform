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
