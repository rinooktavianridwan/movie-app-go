package repositories

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/repository"
    "gorm.io/gorm"
)

type GenreRepository struct {
    DB *gorm.DB
}

func NewGenreRepository(db *gorm.DB) *GenreRepository {
    return &GenreRepository{DB: db}
}

func (r *GenreRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Genre], error) {
    return repository.Paginate[models.Genre](r.DB, page, perPage)
}

func (r *GenreRepository) GetByID(id uint) (*models.Genre, error) {
    var genre models.Genre
    if err := r.DB.First(&genre, id).Error; err != nil {
        return nil, err
    }
    return &genre, nil
}

func (r *GenreRepository) Create(genre *models.Genre) error {
    return r.DB.Create(genre).Error
}

func (r *GenreRepository) Update(genre *models.Genre) error {
    return r.DB.Save(genre).Error
}

func (r *GenreRepository) Delete(id uint) error {
    return r.DB.Delete(&models.Genre{}, id).Error
}

func (r *GenreRepository) ExistsByName(name string) (bool, error) {
    var count int64
    err := r.DB.Model(&models.Genre{}).Where("name = ?", name).Count(&count).Error
    return count > 0, err
}

func (r *GenreRepository) ExistsByNameExceptID(name string, genreID uint) (bool, error) {
    var count int64
    err := r.DB.Model(&models.Genre{}).Where("name = ? AND id != ?", name, genreID).Count(&count).Error
    return count > 0, err
}

func (r *GenreRepository) CountMoviesByGenreID(genreID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.MovieGenre{}).Where("genre_id = ?", genreID).Count(&count).Error
    return count, err
}