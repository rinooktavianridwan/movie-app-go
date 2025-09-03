package services

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/genre/repositories"
	"movie-app-go/internal/modules/genre/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type GenreService struct {
	GenreRepo *repositories.GenreRepository
}

func NewGenreService(genreRepo *repositories.GenreRepository) *GenreService {
	return &GenreService{GenreRepo: genreRepo}
}

func (s *GenreService) CreateGenre(req *requests.CreateGenreRequest) error {
	exists, err := s.GenreRepo.ExistsByName(req.Name)
    if err != nil {
        return err
    }
    if exists {
        return utils.ErrGenreAlreadyExists
    }

	genre := models.Genre{
        Name: req.Name,
    }

    return s.GenreRepo.Create(&genre)
}

func (s *GenreService) GetAllGenresPaginated(page, perPage int) (repository.PaginationResult[models.Genre], error) {
	return s.GenreRepo.GetAllPaginated(page, perPage)
}

func (s *GenreService) GetGenreByID(id uint) (*models.Genre, error) {
	genre, err := s.GenreRepo.GetByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, utils.ErrGenreNotFound
        }
        return nil, err
    }
    return genre, nil
}

func (s *GenreService) UpdateGenre(id uint, req *requests.UpdateGenreRequest) error {
	genre, err := s.GetGenreByID(id)
    if err != nil {
        return err
    }

    if req.Name != genre.Name {
        exists, err := s.GenreRepo.ExistsByNameExceptID(req.Name, id)
        if err != nil {
            return err
        }
        if exists {
            return utils.ErrGenreAlreadyExists
        }
    }

    genre.Name = req.Name

    return s.GenreRepo.Update(genre)
}

func (s *GenreService) DeleteGenre(id uint) error {
	_, err := s.GetGenreByID(id)
    if err != nil {
        return err
    }

    movieCount, err := s.GenreRepo.CountMoviesByGenreID(id)
    if err != nil {
        return err
    }
    if movieCount > 0 {
        return utils.ErrGenreInUse
    }

    return s.GenreRepo.Delete(id)
}
