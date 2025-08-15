package services

import (
	"gorm.io/gorm"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/movie/requests"
)

type MovieService struct {
	DB *gorm.DB
}

func NewMovieService(db *gorm.DB) *MovieService {
	return &MovieService{DB: db}
}

func (s *MovieService) CreateMovie(req *requests.CreateMovieRequest) (*models.Movie, error) {
	var movie models.Movie
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		movie = models.Movie{
			Title:    req.Title,
			Overview: req.Overview,
			Duration: req.Duration,
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
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	if err := s.DB.Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MovieService) GetAllMoviesWithGenres() ([]map[string]interface{}, error) {
	var movies []models.Movie
	if err := s.DB.Find(&movies).Error; err != nil {
		return nil, err
	}

	// Ambil semua relasi movie_genre sekaligus
	var movieGenres []models.MovieGenre
	if err := s.DB.Find(&movieGenres).Error; err != nil {
		return nil, err
	}

	// Mapping movieID -> []genreID
	movieGenreMap := make(map[uint][]uint)
	for _, mg := range movieGenres {
		movieGenreMap[mg.MovieID] = append(movieGenreMap[mg.MovieID], mg.GenreID)
	}

	// Gabungkan data
	var result []map[string]interface{}
	for _, m := range movies {
		result = append(result, map[string]interface{}{
			"id":        m.ID,
			"title":     m.Title,
			"overview":  m.Overview,
			"duration":  m.Duration,
			"genre_ids": movieGenreMap[m.ID],
		})
	}
	return result, nil
}

func (s *MovieService) GetMovieByID(id uint) (*models.Movie, []uint, error) {
	var movie models.Movie
	if err := s.DB.First(&movie, id).Error; err != nil {
		return nil, nil, err
	}
	var movieGenres []models.MovieGenre
	if err := s.DB.Where("movie_id = ?", id).Find(&movieGenres).Error; err != nil {
		return nil, nil, err
	}
	var genreIDs []uint
	for _, mg := range movieGenres {
		genreIDs = append(genreIDs, mg.GenreID)
	}
	return &movie, genreIDs, nil
}

func (s *MovieService) UpdateMovie(id uint, req *requests.UpdateMovieRequest) (*models.Movie, error) {
	var movie models.Movie
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&movie, id).Error; err != nil {
			return err
		}
		movie.Title = req.Title
		movie.Overview = req.Overview
		movie.Duration = req.Duration
		if err := tx.Save(&movie).Error; err != nil {
			return err
		}
		// Hapus relasi lama
		if err := tx.Where("movie_id = ?", id).Delete(&models.MovieGenre{}).Error; err != nil {
			return err
		}
		// Tambah relasi baru
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
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *MovieService) DeleteMovie(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("movie_id = ?", id).Delete(&models.MovieGenre{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Movie{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
