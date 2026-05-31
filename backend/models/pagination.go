package models

type PaginationResponse struct {
	Items       interface{} `json:"items"`
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
}

// CursorPaginationResponse 游标分页响应
type CursorPaginationResponse struct {
	Items      interface{} `json:"items"`
	NextCursor string      `json:"next_cursor,omitempty"`
	PrevCursor string      `json:"prev_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
	Total      int64       `json:"total"`
}

// CursorPageParams 游标分页参数
type CursorPageParams struct {
	Cursor   string // 上一次最后一条记录的ID
	Limit    int    // 每页大小，默认20
	Direction string // 方向: "next" 或 "prev"
}