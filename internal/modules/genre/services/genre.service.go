package services

import (
	"fmt"
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"

	"gorm.io/gorm"
)

type GenreService struct {
	DB *gorm.DB
}

func NewGenreService(db *gorm.DB) *GenreService {
	return &GenreService{DB: db}
}

func (s *GenreService) CreateGenre(name string) (*models.Genre, error) {
	genre := models.Genre{Name: name}
	if err := s.DB.Create(&genre).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

func (s *GenreService) GetAllGenresPaginated(page, perPage int) (repository.PaginationResult[models.Genre], error) {
	return repository.Paginate[models.Genre](s.DB, page, perPage)
}

func (s *GenreService) GetGenreByID(id uint) (*models.Genre, error) {
	var genre models.Genre
	if err := s.DB.First(&genre, id).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

func (s *GenreService) UpdateGenre(id uint, name string) (*models.Genre, error) {
	var genre models.Genre
	if err := s.DB.First(&genre, id).Error; err != nil {
		return nil, err
	}
	genre.Name = name
	if err := s.DB.Save(&genre).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

func (s *GenreService) DeleteGenre(id uint) error {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        var genre models.Genre
        if err := tx.First(&genre, id).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                return fmt.Errorf("genre not found")
            }
            return err
        }
        
        var count int64
        if err := tx.Model(&models.MovieGenre{}).Where("genre_id = ?", id).Count(&count).Error; err != nil {
            return err
        }
        if count > 0 {
            return fmt.Errorf("cannot delete genre: still used by movies")
        }

        if err := tx.Delete(&models.Genre{}, id).Error; err != nil {
            return err
        }
        return nil
    })
}
