package repositories

import "github.com/google/uuid"

type PageFilter struct {
	Page     int
	PageSize int
}

func (f PageFilter) LimitOffset() (int, int) {
	page := f.Page
	if page < 1 {
		page = 1
	}
	size := f.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return size, (page - 1) * size
}

type ListResult[T any] struct {
	Items []T
	Total int
}

type IDFilter struct {
	ID uuid.UUID
}
