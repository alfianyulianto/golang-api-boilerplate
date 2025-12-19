package response

import "math"

type Response[T any] struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       T           `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Error      any         `json:"error,omitempty"`
}
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItem   int64 `json:"total_item"`
	TotalPage   int   `json:"total_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

type PaginationRequest struct {
	Page     int `json:"page" validate:"omitempty,min=1"`
	PageSize int `json:"page_size" validate:"omitempty,min=1,max=100"`
}

func ToPaginated(page, pageSize int, totalItem int64) *Pagination {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	totalPage := int(math.Ceil(float64(totalItem) / float64(pageSize)))

	return &Pagination{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItem:   totalItem,
		TotalPage:   totalPage,
		HasNext:     page < totalPage,
		HasPrev:     page > 1,
	}
}
