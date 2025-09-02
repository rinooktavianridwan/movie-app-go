package services

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type GenreService struct {
	DB *gorm.DB
}

func NewGenreService(db *gorm.DB) *GenreService {
	return &GenreService{DB: db}
}

func (s *GenreService) CreateGenre(name string) error {
	genre := models.Genre{Name: name}
	if err := s.DB.Create(&genre).Error; err != nil {
		return err
	}
	return nil
}

func (s *GenreService) GetAllGenresPaginated(page, perPage int) (repository.PaginationResult[models.Genre], error) {
	return repository.Paginate[models.Genre](s.DB, page, perPage)
}

func (s *GenreService) GetGenreByID(id uint) (*models.Genre, error) {
	var genre models.Genre
	if err := s.DB.First(&genre, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrGenreNotFound
		}
		return nil, err
	}
	return &genre, nil
}

func (s *GenreService) UpdateGenre(id uint, name string) error {
	var genre models.Genre
	if err := s.DB.First(&genre, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrGenreNotFound
		}
		return err
	}

	genre.Name = name
	if err := s.DB.Save(&genre).Error; err != nil {
		return err
	}
	return nil
}

func (s *GenreService) DeleteGenre(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var genre models.Genre
		if err := tx.First(&genre, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrGenreNotFound
			}
			return err
		}

		var movieGenreCount int64
		if err := tx.Model(&models.MovieGenre{}).Where("genre_id = ?", id).Count(&movieGenreCount).Error; err != nil {
			return err
		}
		if movieGenreCount > 0 {
			return utils.ErrGenreInUse
		}

		if err := tx.Delete(&models.Genre{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
