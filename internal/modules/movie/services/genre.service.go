package services

import (
    "gorm.io/gorm"
    "movie-app-go/internal/models"
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

func (s *GenreService) GetAllGenres() ([]models.Genre, error) {
    var genres []models.Genre
    if err := s.DB.Find(&genres).Error; err != nil {
        return nil, err
    }
    return genres, nil
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
    return s.DB.Delete(&models.Genre{}, id).Error
}