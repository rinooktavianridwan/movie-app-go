package services

import (
	"mime/multipart"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/movie/repositories"
	"movie-app-go/internal/modules/movie/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"strings"

	"gorm.io/gorm"
)

type MovieService struct {
	MovieRepo *repositories.MovieRepository
}

func NewMovieService(movieRepo *repositories.MovieRepository) *MovieService {
	return &MovieService{
		MovieRepo: movieRepo,
	}
}

func (s *MovieService) CreateMovie(req *requests.CreateMovieRequest, posterFile *multipart.FileHeader) error {
	count, err := s.MovieRepo.CountGenresByIDs(req.GenreIDs)
	if err != nil {
		return err
	}
	if count != int64(len(req.GenreIDs)) {
		return utils.ErrInvalidGenreIDs
	}

	var posterURL *string
	if posterFile != nil {
		posterPath, err := utils.SaveFile(posterFile, "uploads/posters", "image", 10)
		if err != nil {
			return err
		}
		relativePath := strings.TrimPrefix(posterPath, "./")
		posterURL = &relativePath
	}

	movie := models.Movie{
		Title:     req.Title,
		Overview:  req.Overview,
		Duration:  req.Duration,
		PosterURL: posterURL,
	}

	return s.MovieRepo.WithTransaction(func(tx *gorm.DB) error {
		if err := s.MovieRepo.CreateWithTx(tx, &movie); err != nil {
			return err
		}

		movieGenres := make([]models.MovieGenre, 0, len(req.GenreIDs))
		for _, gid := range req.GenreIDs {
			movieGenres = append(movieGenres, models.MovieGenre{
				MovieID: movie.ID,
				GenreID: gid,
			})
		}

		return s.MovieRepo.CreateMovieGenresWithTx(tx, movieGenres)
	})
}

func (s *MovieService) GetAllMoviesPaginated(page, perPage int) (repository.PaginationResult[models.Movie], error) {
	return s.MovieRepo.GetAllPaginated(page, perPage)
}

func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	movie, err := s.MovieRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) UpdateMovie(id uint, req *requests.UpdateMovieRequest, posterFile *multipart.FileHeader) error {
	movie, err := s.GetMovieByID(id)
	if err != nil {
		return err
	}

	count, err := s.MovieRepo.CountGenresByIDs(req.GenreIDs)
	if err != nil {
		return err
	}
	if count != int64(len(req.GenreIDs)) {
		return utils.ErrInvalidGenreIDs
	}

	if posterFile != nil {
		if movie.PosterURL != nil && *movie.PosterURL != "" {
			utils.DeleteFile(*movie.PosterURL)
		}

		posterPath, err := utils.SaveFile(posterFile, "uploads/posters", "image", 10)
		if err != nil {
			return err
		}
		relativePath := strings.TrimPrefix(posterPath, "./")
		movie.PosterURL = &relativePath
	}

	movie.Title = req.Title
	movie.Overview = req.Overview
	movie.Duration = req.Duration

	return s.MovieRepo.WithTransaction(func(tx *gorm.DB) error {
		if err := s.MovieRepo.UpdateWithTx(tx, movie); err != nil {
			return err
		}

		if err := s.MovieRepo.DeleteMovieGenresWithTx(tx, movie.ID); err != nil {
			return err
		}

		movieGenres := make([]models.MovieGenre, 0, len(req.GenreIDs))
		for _, gid := range req.GenreIDs {
			movieGenres = append(movieGenres, models.MovieGenre{
				MovieID: movie.ID,
				GenreID: gid,
			})
		}

		return s.MovieRepo.CreateMovieGenresWithTx(tx, movieGenres)
	})
}

func (s *MovieService) DeleteMovie(id uint) error {
	movie, err := s.GetMovieByID(id)
	if err != nil {
		return err
	}

	return s.MovieRepo.WithTransaction(func(tx *gorm.DB) error {
		scheduleCount, err := s.MovieRepo.CountSchedulesByMovieID(id)
		if err != nil {
			return err
		}
		if scheduleCount > 0 {
			return utils.ErrMovieHasSchedules
		}

		if err := s.MovieRepo.DeleteMovieGenresWithTx(tx, id); err != nil {
			return err
		}

		if movie.PosterURL != nil && *movie.PosterURL != "" {
			utils.DeleteFile(*movie.PosterURL)
		}

		if err := s.MovieRepo.DeleteWithTx(tx, id); err != nil {
			return err
		}

		return nil
	})
}
