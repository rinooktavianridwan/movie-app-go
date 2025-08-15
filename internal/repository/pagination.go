package repository

import (
	"gorm.io/gorm"
)

type PaginationResult[T any] struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
    Data       []T   `json:"data"`
}

func Paginate[T any](db *gorm.DB, page, perPage int) (PaginationResult[T], error) {
	var result PaginationResult[T]
	var total int64
	var data []T

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	db.Model(new(T)).Count(&total)
	offset := (page - 1) * perPage

	if err := db.Limit(perPage).Offset(offset).Find(&data).Error; err != nil {
		return result, err
	}

	result = PaginationResult[T]{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
		Data:       data,
	}
	return result, nil
}
