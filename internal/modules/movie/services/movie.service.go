package services

import (
	"mime/multipart"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/movie/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"strings"

	"gorm.io/gorm"
)

type MovieService struct {
	DB *gorm.DB
}

func NewMovieService(db *gorm.DB) *MovieService {
	return &MovieService{DB: db}
}

func (s *MovieService) CreateMovie(req *requests.CreateMovieRequest, posterFile *multipart.FileHeader) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Genre{}).Where("id IN ?", req.GenreIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(req.GenreIDs)) {
			return utils.ErrInvalidGenreIDs
		}

		movie := models.Movie{
			Title:    req.Title,
			Overview: req.Overview,
			Duration: req.Duration,
		}

		if posterFile != nil {
			posterPath, err := utils.SaveFile(posterFile, "uploads/posters", "image", 10)
			if err != nil {
				return err
			}
			relativePath := strings.TrimPrefix(posterPath, "./")
			movie.PosterURL = &relativePath
		}

		if err := tx.Create(&movie).Error; err != nil {
			return err
		}

		movieGenres := make([]models.MovieGenre, 0, len(req.GenreIDs))
		for _, gid := range req.GenreIDs {
			movieGenres = append(movieGenres, models.MovieGenre{
				MovieID: movie.ID,
				GenreID: gid,
			})
		}
		if len(movieGenres) > 0 {
			if err := tx.Create(&movieGenres).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *MovieService) GetAllMoviesPaginated(page, perPage int) (repository.PaginationResult[models.Movie], error) {
	return repository.Paginate[models.Movie](
		s.DB.Preload("MovieGenres.Genre"),
		page,
		perPage,
	)
}

func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := s.DB.Preload("MovieGenres.Genre").First(&movie, id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *MovieService) UpdateMovie(id uint, req *requests.UpdateMovieRequest, posterFile *multipart.FileHeader) error {
	var movie models.Movie
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&movie, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrMovieNotFound
			}
			return err
		}

		var count int64
		if err := tx.Model(&models.Genre{}).Where("id IN ?", req.GenreIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(req.GenreIDs)) {
			return utils.ErrInvalidGenreIDs
		}

		movie.Title = req.Title
		movie.Overview = req.Overview
		movie.Duration = req.Duration

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

		if err := tx.Save(&movie).Error; err != nil {
			return err
		}

		if err := tx.Where("movie_id = ?", id).Delete(&models.MovieGenre{}).Error; err != nil {
			return err
		}

		movieGenres := make([]models.MovieGenre, 0, len(req.GenreIDs))
		for _, gid := range req.GenreIDs {
			movieGenres = append(movieGenres, models.MovieGenre{
				MovieID: movie.ID,
				GenreID: gid,
			})
		}
		if len(movieGenres) > 0 {
			if err := tx.Create(&movieGenres).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *MovieService) DeleteMovie(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var movie models.Movie
		if err := tx.First(&movie, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrMovieNotFound
			}
			return err
		}

		var scheduleCount int64
		if err := tx.Model(&models.Schedule{}).Where("movie_id = ?", id).Count(&scheduleCount).Error; err != nil {
			return err
		}
		if scheduleCount > 0 {
			return utils.ErrMovieHasSchedules
		}

		if err := tx.Where("movie_id = ?", id).Delete(&models.MovieGenre{}).Error; err != nil {
            return err
        }

		if movie.PosterURL != nil && *movie.PosterURL != "" {
			utils.DeleteFile(*movie.PosterURL)
		}

		if err := tx.Delete(&models.Movie{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
