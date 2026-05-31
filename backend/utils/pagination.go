package utils

import (
	"warehouse-management/backend/models"

	"gorm.io/gorm"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func GetCursorParams(cursor string, limit int, direction string) models.CursorPageParams {
	if limit <= 0 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if direction == "" {
		direction = "next"
	}

	return models.CursorPageParams{
		Cursor:    cursor,
		Limit:     limit,
		Direction: direction,
	}
}

func ApplyCursorPagination(db *gorm.DB, params models.CursorPageParams, IDField string) *gorm.DB {
	if params.Cursor == "" {
		return db.Limit(params.Limit)
	}

	if params.Direction == "prev" {
		return db.Where(IDField+" < ?", params.Cursor).Order(IDField + " DESC").Limit(params.Limit)
	}

	return db.Where(IDField+" > ?", params.Cursor).Order(IDField + " ASC").Limit(params.Limit)
}

func BuildCursorResponse(items interface{}, lastID string, hasMore bool, total int64) models.CursorPaginationResponse {
	var nextCursor, prevCursor string

	if lastID != "" {
		if hasMore {
			nextCursor = lastID
		}
		prevCursor = lastID
	}

	return models.CursorPaginationResponse{
		Items:      items,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasMore:    hasMore,
		Total:      total,
	}
}