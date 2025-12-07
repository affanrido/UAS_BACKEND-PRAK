package model

// PaginationRequest - DTO untuk pagination request
type PaginationRequest struct {
	Page     int `json:"page"`      // Current page (1-based)
	PageSize int `json:"page_size"` // Items per page
}

// PaginationResponse - DTO untuk pagination response
type PaginationResponse struct {
	Page       int `json:"page"`        // Current page
	PageSize   int `json:"page_size"`   // Items per page
	TotalItems int `json:"total_items"` // Total items
	TotalPages int `json:"total_pages"` // Total pages
}

// GetOffset - Calculate offset for SQL query
func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit - Get limit for SQL query
func (p *PaginationRequest) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10 // Default page size
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Max page size
	}
	return p.PageSize
}

// CalculateTotalPages - Calculate total pages
func CalculateTotalPages(totalItems, pageSize int) int {
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
