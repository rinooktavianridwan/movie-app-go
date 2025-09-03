package repositories

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"

	"gorm.io/gorm"
)

type MovieRepository struct {
	DB *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		DB: db,
	}
}

func (r *MovieRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Movie], error) {
	return repository.Paginate[models.Movie](
		r.DB.Preload("MovieGenres.Genre"),
		page,
		perPage,
	)
}

func (r *MovieRepository) GetByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := r.DB.Preload("MovieGenres.Genre").First(&movie, id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) CountGenresByIDs(genreIDs []uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Genre{}).Where("id IN ?", genreIDs).Count(&count).Error
	return count, err
}

func (r *MovieRepository) CountSchedulesByMovieID(movieID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Schedule{}).Where("movie_id = ?", movieID).Count(&count).Error
	return count, err
}

// Transaction Methods:
func (r *MovieRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}

func (r *MovieRepository) CreateWithTx(tx *gorm.DB, movie *models.Movie) error {
	return tx.Create(movie).Error
}

func (r *MovieRepository) CreateMovieGenresWithTx(tx *gorm.DB, movieGenres []models.MovieGenre) error {
	if len(movieGenres) == 0 {
		return nil
	}
	return tx.Create(&movieGenres).Error
}

func (r *MovieRepository) UpdateWithTx(tx *gorm.DB, movie *models.Movie) error {
	return tx.Save(movie).Error
}

func (r *MovieRepository) DeleteWithTx(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Movie{}, id).Error
}

func (r *MovieRepository) DeleteMovieGenresWithTx(tx *gorm.DB, movieID uint) error {
	return tx.Where("movie_id = ?", movieID).Delete(&models.MovieGenre{}).Error
}
