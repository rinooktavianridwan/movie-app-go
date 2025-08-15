package repository

import "gorm.io/gorm"

func PaginateWithPreload[T any](db *gorm.DB, page, perPage int, preloads ...string) (PaginationResult[T], error) {
    var result PaginationResult[T]
    var total int64
    var data []T

    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = 10
    }

    tx := db.Model(new(T))
    tx.Count(&total)
    offset := (page - 1) * perPage

    // Apply preloads
    for _, preload := range preloads {
        tx = tx.Preload(preload)
    }

    if err := tx.Limit(perPage).Offset(offset).Find(&data).Error; err != nil {
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