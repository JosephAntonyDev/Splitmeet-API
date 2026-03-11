package core

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page    int
	Limit   int
	Offset  int
	Search  string
	SortBy  string
	SortDir string
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Sanitize sort direction
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return PaginationParams{
		Page:    page,
		Limit:   limit,
		Offset:  (page - 1) * limit,
		Search:  search,
		SortBy:  sortBy,
		SortDir: sortDir,
	}
}

func NewPaginatedResponse(data interface{}, page, limit, total int) PaginatedResponse {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// SanitizeSortColumn validates that the sort column is in the allowed list to prevent SQL injection
func SanitizeSortColumn(column string, allowed []string) string {
	for _, a := range allowed {
		if column == a {
			return column
		}
	}
	return "created_at"
}

// BuildOrderClause creates a safe ORDER BY clause
func BuildOrderClause(sortBy, sortDir string, allowedColumns []string) string {
	col := SanitizeSortColumn(sortBy, allowedColumns)
	if sortDir != "asc" {
		sortDir = "desc"
	}
	return fmt.Sprintf("%s %s", col, sortDir)
}
